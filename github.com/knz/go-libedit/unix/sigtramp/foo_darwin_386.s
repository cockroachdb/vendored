#include "go_asm.h"
#include "textflag.h"

TEXT ·getptr(SB),NOSPLIT,$0
	MOVL    $runtime·sigtramp(SB), AX
	MOVL    AX, ret+0(FP)
	RET
