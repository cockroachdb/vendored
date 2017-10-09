#include <histedit.h>
#include <stdint.h>
#include <errno.h>
#include <stdio.h>
#include <setjmp.h>
#include <termios.h>
#include <signal.h>

#include "c_editline.h"

char *go_libedit_emptycstring = (char*)"";
const char* go_libedit_mode_read = "r";
const char* go_libedit_mode_write = "w";
const char* go_libedit_mode_append = "a";
const char* go_libedit_locale1 = "en_US.UTF-8";
const char* go_libedit_locale2 = "C.UTF-8";

go_libedit_promptgen go_libedit_prompt_left_ptr = go_libedit_prompt_left;
go_libedit_promptgen go_libedit_prompt_right_ptr = go_libedit_prompt_right;

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

void go_libedit_set_prompt(EditLine *el, int p, go_libedit_promptgen f) {
    el_set(el, p, f);
}

static unsigned char	 _el_rl_complete(EditLine *, int);
static unsigned char	 _el_rl_tstp(EditLine *, int);
static unsigned char	 _el_rl_intr(EditLine *, int);

void go_libedit_rebind_ctrls(EditLine *e) {
    // Handle ^C properly.
    el_set(e, EL_ADDFN, "rl_interrupt",
	   "ReadLine compatible interrupt function",
	   _el_rl_intr);
    el_set(e, EL_BIND, "^C", "rl_interrupt", NULL);

    // Word completion - this has to go after loading the default
    // mappings.
    el_set(e, EL_ADDFN, "rl_complete",
	   "ReadLine compatible completion function",
	   _el_rl_complete);
    el_set(e, EL_BIND, "^I", "rl_complete", NULL);

    // Send TSTP when ^Z is pressed.
    el_set(e, EL_ADDFN, "rl_tstp",
	   "ReadLine compatible suspend function",
	   _el_rl_tstp);
    el_set(e, EL_BIND, "^Z", "rl_tstp", NULL);

    // Readline history search. People are used to this.
    el_set(e, EL_BIND, "^R", "em-inc-search-prev", NULL);
}

EditLine* go_libedit_init(char *appName,
			  FILE* fin, FILE* fout, FILE *ferr) {
    // Create the editor.
    EditLine *e = el_init(appName, fin, fout, ferr);
    if (!e) {
	return NULL;
    }

    // Do we really want to edit?
    int editmode = 1;
    struct termios t;
    if (tcgetattr(fileno(fin), &t) != -1 && (t.c_lflag & ECHO) == 0)
	editmode = 0;
    if (!editmode)
	el_set(e, EL_EDITMODE, 0);

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

static unsigned char _el_rl_tstp(EditLine *el, int ch) {
    (void) kill(0, SIGTSTP);
    return CC_NORM;
}

static sigjmp_buf jmpbuf;
static unsigned char _el_rl_intr(EditLine *el, int ch) {
    // Reveal the Ctrl+C to the top-level caller of el_gets.
    siglongjmp(jmpbuf, 1);
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

void *go_libedit_gets(EditLine *el, int *count, int *interrupted, int widechar) {
    void *ret = NULL;
    int saveerr = 0;

    // Disable conversion of Ctrl+C to signal.
    FILE *inf;
    el_get(el, EL_GETFP, 0, &inf);
    struct termios t;
    int intr_disabled = 0;
    cc_t intr_char;
    if (tcgetattr(fileno(inf), &t) != -1) {
	intr_char = t.c_cc[VINTR];
	t.c_cc[VINTR] = 0;
	tcsetattr(fileno(inf), TCSANOW, &t);
	intr_disabled = 1;
    }

    // Set up libedit's signal handlers.
    // We need to do this on every invocaiton of el_gets
    // because cgo resets signal handlers at the C/Go boundary.
    el_set(el, EL_SIGNAL, 1);

    // Prepare to be interrupted.
    // This will occur when Ctrl+C is entered at the beginning of a
    // line.
    if (sigsetjmp(jmpbuf, 1)) {
	saveerr = EINTR;
	*interrupted = 1;
	ret = NULL;
	goto restore;
    }

    // Read the line.
    if (widechar) {
	ret = (void *)el_wgets(el, count);
    } else {
	ret = (void *)el_gets(el, count);
    }
    saveerr = errno;

restore:
    // Remove libedit's signal handlers.
    el_set(el, EL_SIGNAL, 0);

    // Restore Ctrl+C processing by the terminal.
    if (intr_disabled) {
	t.c_cc[VINTR] = intr_char;
	tcsetattr(fileno(inf), TCSANOW, &t);
    }

    // Restore errno.
    errno = saveerr;
    return ret;
}
