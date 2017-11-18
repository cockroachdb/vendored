#include "go_asm.h"
#include "textflag.h"

TEXT ·getptr(SB),NOSPLIT,$0
	MOVW    $runtime·sigtramp(SB), R0
	MOVW    R0, ret+0(FP)
	RET
