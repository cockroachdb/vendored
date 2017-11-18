#ifndef GO_LIBEDIT_H
#define GO_LIBEDIT_H

#include <histedit.h>
#include <stdio.h>

typedef char* pchar;

EditLine* go_libedit_init(char *appName, void **sigcfg,
			  FILE* fin, FILE* fout, FILE *ferr,
			  void *sigtramp);
void go_libedit_close(EditLine *el, void *sigcfg);
void go_libedit_rebind_ctrls(EditLine *el);

extern char *go_libedit_emptycstring;
extern const char* go_libedit_mode_read;
extern const char* go_libedit_mode_write;
extern const char* go_libedit_mode_append;
extern const char *go_libedit_locale1;
extern const char *go_libedit_locale2;

int go_libedit_get_clientdata(EditLine *el);
void go_libedit_set_clientdata(EditLine *el, int v);
void go_libedit_set_string_array(char **ar, int p, char *s);

void *go_libedit_gets(EditLine *el, char *lprompt, char *rprompt,
		      void *sigcfg, int *count, int *interrupted, int wc);

History* go_libedit_setup_history(EditLine *el, int maxEntries, int dedup);
int go_libedit_read_history(History *h, char *filename);
int go_libedit_write_history(History *h, char *filename);
int go_libedit_add_history(History *h, char *line);

// Go-generated via //export
char *go_libedit_prompt_left(EditLine *el);
char *go_libedit_prompt_right(EditLine *el);
char **go_libedit_getcompletions(int instance, char *word);


#endif
