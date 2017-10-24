#include <histedit.h>
#include <stdint.h>
#include <errno.h>
#include <stdio.h>
#include <setjmp.h>
#include <termios.h>
#include <signal.h>
#include <stdlib.h>
#include <unistd.h>

#include "c_editline.h"

char *go_libedit_emptycstring = (char*)"";
const char* go_libedit_mode_read = "r";
const char* go_libedit_mode_write = "w";
const char* go_libedit_mode_append = "a";
const char* go_libedit_locale1 = "en_US.UTF-8";
const char* go_libedit_locale2 = "C.UTF-8";

void go_libedit_set_string_array(char **ar, int p, char *s) {
    ar[p] = s;
}

int go_libedit_get_clientdata(EditLine *el) {
    void* p;
    el_get(el, EL_CLIENTDATA, &p);
    return (int)(intptr_t)p;
}

void go_libedit_set_clientdata(EditLine *el, int v) {
    void *p = (void*)(intptr_t) v;
    el_set(el, EL_CLIENTDATA, p);
}

/************** prompts **************/

static char *go_libedit_lprompt_str;
static char *go_libedit_rprompt_str;

static char* go_libedit_lprompt(EditLine *el) {
    return go_libedit_lprompt_str ? go_libedit_lprompt_str : "?";
}

static char* go_libedit_rprompt(EditLine *el) {
    return go_libedit_rprompt_str ? go_libedit_rprompt_str : "";
}


/********** initialization ***********/

static unsigned char	 _el_rl_complete(EditLine *, int);

static volatile int g_interrupted;
static unsigned char	 _el_rl_intr(EditLine * e, int c) {
    g_interrupted = 1;
    return CC_NEWLINE;
}

static unsigned char	 _el_rl_tstp(EditLine * e, int c) {
    kill(0, SIGTSTP);
    return CC_REDISPLAY;
}


void go_libedit_rebind_ctrls(EditLine *e) {
    // Handle Ctrl+C gracefully.
    el_set(e, EL_ADDFN, "rl_interrupt",
	  "ReadLine compatible interrupt function",
	  _el_rl_intr);
    el_set(e, EL_BIND, "^C", "rl_interrupt", NULL);

    // Handle Ctrl+Z gracefully.
    el_set(e, EL_ADDFN, "rl_tstp",
	   "ReadLine compatible suspend function",
	   _el_rl_tstp);
    el_set(e, EL_BIND, "^Z", "rl_tstp", NULL);

    // Word completion - this has to go after loading the default
    // mappings.
    el_set(e, EL_ADDFN, "rl_complete",
	   "ReadLine compatible completion function",
	   _el_rl_complete);
    el_set(e, EL_BIND, "^I", "rl_complete", NULL);

    // Readline history search. People are used to this.
    el_set(e, EL_BIND, "^R", "em-inc-search-prev", NULL);
}

EditLine* go_libedit_init(char *appName, void** el_signal,
			  FILE* fin, FILE* fout, FILE *ferr) {
    // Prepare signal handling.
    (*el_signal) = 0;

    // Do we really want to edit?
    int editmode = 1;
    struct termios t;
    if (tcgetattr(fileno(fin), &t) != -1 && (t.c_lflag & ECHO) == 0)
	editmode = 0;

    // Create the editor.

    // We need to mark here in the terminal config that Ctrl+C doesn't
    // do anything, so that el_init() captures and saves that
    // information. It will be used every time el_gets() starts,
    // when reset the terminal to "known settings".
    EditLine *e = el_init(appName, fin, fout, ferr);

    if (!e) {
	return NULL;
    }

    if (!editmode)
	el_set(e, EL_EDITMODE, 0);

    // Set up the prompt functions. Unfortunately we cannot use a real
    // callback into Go because Go doesn't like being called from an
    // alternate signal stack.
    el_set(e, EL_PROMPT, go_libedit_lprompt);
    el_set(e, EL_RPROMPT, go_libedit_rprompt);

    // Load the emacs keybindings by default. We need
    // to do that before the defaults are overridden below.
    el_set(e, EL_EDITOR, "emacs");

    go_libedit_rebind_ctrls(e);

    // Home/End keys.
    el_set(e, EL_BIND, "\\e[1~", "ed-move-to-beg", NULL);
    el_set(e, EL_BIND, "\\e[4~", "ed-move-to-end", NULL);
    el_set(e, EL_BIND, "\\e[7~", "ed-move-to-beg", NULL);
    el_set(e, EL_BIND, "\\e[8~", "ed-move-to-end", NULL);
    el_set(e, EL_BIND, "\\e[H", "ed-move-to-beg", NULL);
    el_set(e, EL_BIND, "\\e[F", "ed-move-to-end", NULL);

    // Delete/Insert keys.
    el_set(e, EL_BIND, "\\e[3~", "ed-delete-next-char", NULL);
    el_set(e, EL_BIND, "\\e[2~", "ed-quoted-insert", NULL);

    // Ctrl-left-arrow and Ctrl-right-arrow for word moving.
    el_set(e, EL_BIND, "\\e[1;5C", "em-next-word", NULL);
    el_set(e, EL_BIND, "\\e[1;5D", "ed-prev-word", NULL);
    el_set(e, EL_BIND, "\\e[5C", "em-next-word", NULL);
    el_set(e, EL_BIND, "\\e[5D", "ed-prev-word", NULL);
    el_set(e, EL_BIND, "\\e\\e[C", "em-next-word", NULL);
    el_set(e, EL_BIND, "\\e\\e[D", "ed-prev-word", NULL);

    // Read the settings from the configuration file.
    el_source(e, NULL);

    return e;
}

/************** history **************/

History* go_libedit_setup_history(EditLine *el, int maxEntries, int dedup) {
    if (!el) {
	errno = EINVAL;
	return NULL;
    }

    History *h = history_init();
    if (!h)
	return NULL;

    HistEvent ev;
    history(h, &ev, H_SETSIZE, maxEntries);
    history(h, &ev, H_SETUNIQUE, dedup);

    el_set(el, EL_HIST, history, h);
    return h;
}

static int readwrite_history(History *h, int op, char *filename) {
    if (!h || !filename) {
	errno = EINVAL;
	return -1;
    }
    errno = 0;
    HistEvent ev;
    int res;
    if ((res = history(h, &ev, op, filename)) == -1) {
	if (!errno)
	    errno = EINVAL;
	return -1;
    }
    return res;
}

int go_libedit_read_history(History *h, char *filename) {
    return readwrite_history(h, H_LOAD, filename);
}

int go_libedit_write_history(History *h, char *filename) {
    return readwrite_history(h, H_SAVE, filename);
}

int go_libedit_add_history(History *h, char *line) {
    return readwrite_history(h, H_ENTER, line);
}


/************* completion ************/

// We can't use rl_complete directly because that uses the readline
// emulation's own EditLine instance, and here we want to use our
// own. So basically re-implement on top of editline's internal
// fn_complete function.

int
fn_complete(EditLine *el,
	    char *(*complet_func)(const char *, int),
	    char **(*attempted_completion_function)(const char *, int, int),
	    const wchar_t *word_break, const wchar_t *special_prefixes,
	    const char *(*app_func)(const char *), size_t query_items,
	    int *completion_type, int *over, int *point, int *end,
	    const wchar_t *(*find_word_start_func)(const wchar_t *, const wchar_t *),
	    wchar_t *(*dequoting_func)(const wchar_t *),
	    char *(*quoting_func)(const char *));
static const wchar_t break_chars[] = L" \t\n\"\\'`@$><=;|&{(";


// In an unfortunate turn of circumstances, editline's fn_complete
// API does not pass the EditLine instance nor the clientdata field
// to the attempted_completion_function, yet we really want this.
// So we'll pass it as a hidden argument via a global variable.
// This effectively makes the entire library thread-unsafe. :'-(

static int global_instance;

static char **wrap_autocomplete(const char *word, int unused1, int unused2) {
    return go_libedit_getcompletions(global_instance, (char*)word);
}

static const char *_rl_completion_append_character_function(const char *_) {
    static const char *sp = " ";
    return sp;
}

static unsigned char _el_rl_complete(EditLine *el, int ch) {
    int avoid_filename_completion = 1;

    // Urgh...
    global_instance = go_libedit_get_clientdata(el);

    return (unsigned char)fn_complete(
	el,
	NULL /* complet_func */,
	wrap_autocomplete /* attempted_completion_function */,
	break_chars /* word_break */,
	NULL /* special_prefixes */,
	_rl_completion_append_character_function /* app_func */,
	100 /* query_items */,
	NULL /* completion_type */,
	&avoid_filename_completion /* over */,
	NULL /* point */,
	NULL /* end */,
	NULL /* find_word_start_func */,
	NULL /* dequoting_func */,
	NULL /* quoting_func */
	);
}


/*************** el_gets *************/

static EditLine* volatile g_el;
static volatile int g_fdout;

static void winch_handler(int signo) {
    // Tell libedit about the window size change.
    el_resize(g_el);

    // Now we want to redraw the current line, however libedit at this
    // point expects the cursor to be at the beginning of the
    // line. Make it so.
    write(g_fdout, "\r", 1);

    // Ready to refresh. Do it.
    el_set(g_el, EL_REFRESH);
}

static void cont_handler(int signo) {
    // This can happen when the process comes back to
    // the foreground, after being suspended and
    // perhaps after a shell has reset the terminal.
    //
    // So refresh the input stream's termios.
    // We don't really need to reset intr here, just
    // to call SETTY -d so that the editor's termios
    // get re-loaded.
    el_set(g_el, EL_SETTY, "-d", "intr=", NULL);

    // Then re-load the terminal settings (size etc).
    el_set(g_el, EL_PREP_TERM, 1);
    // And redraw the current input line.
    el_set(g_el, EL_REFRESH);
}


void *go_libedit_gets(EditLine *el, char *lprompt, char *rprompt,
		      void *p_sigcfg, int *count, int *interrupted, int widechar) {
    void *ret = NULL;
    int saveerr = 0;

    // Save and clear Go's signal handlers, set up our own for SIGCONT and SIGWINCH.
    struct sigaction osa[2];
    struct sigaction nsa;
    g_el = el;
    FILE *fout;
    el_get(el, EL_GETFP, 1, &fout);
    g_fdout = fileno(fout);
    sigaction(SIGCONT, 0, &nsa);
    nsa.sa_handler = cont_handler;
    nsa.sa_flags |= SA_RESTART;
    sigaction(SIGCONT, &nsa, &osa[0]);
    sigaction(SIGWINCH, 0, &nsa);
    nsa.sa_handler = winch_handler;
    nsa.sa_flags |= SA_RESTART;
    sigaction(SIGWINCH, &nsa, &osa[1]);

    // Disable terminal translation of ^C and ^Z to SIGINT and SIGTSTP.
    el_set(el, EL_SETTY, "-d", "intr=", NULL);
    el_set(el, EL_SETTY, "-d", "susp=", NULL);
    el_set(el, EL_SETTY, "-x", "intr=", NULL);
    el_set(el, EL_SETTY, "-x", "susp=", NULL);
    el_set(el, EL_SETTY, "-q", "intr=", NULL);
    el_set(el, EL_SETTY, "-q", "susp=", NULL);
    el_reset(el);

    // Install the prompts.
    go_libedit_lprompt_str = lprompt;
    go_libedit_rprompt_str = rprompt;

    // Read the line.
    g_interrupted = 0;
    if (widechar) {
	ret = (void *)el_wgets(el, count);
    } else {
	ret = (void *)el_gets(el, count);
    }
    // Save errno so we have something meaningful to return below.
    saveerr = errno;
    if (g_interrupted != 0) {
	saveerr = EINTR;
	*interrupted = 1;
	ret = NULL;
    }

    go_libedit_lprompt_str = NULL;
    go_libedit_rprompt_str = NULL;

    // Restore Go's signal handlers.
    sigaction(SIGCONT, &osa[0], 0);
    sigaction(SIGWINCH, &osa[1], 0);

    // If libedit got interrupted it may not have restored the
    // terminal settings properly. Restore them.
    el_set(el, EL_SETTY, "-d", "intr=^C", NULL);
    el_set(el, EL_SETTY, "-d", "susp=^Z", NULL);
    el_set(el, EL_SETTY, "-x", "intr=^C", NULL);
    el_set(el, EL_SETTY, "-x", "susp=^Z", NULL);
    el_set(el, EL_SETTY, "-q", "intr=^C", NULL);
    el_set(el, EL_SETTY, "-q", "susp=^Z", NULL);

    // Restore errno.
    errno = saveerr;
    return ret;
}

void go_libedit_close(EditLine *el, void *p_sigcfg) {
    el_end(el);
}
