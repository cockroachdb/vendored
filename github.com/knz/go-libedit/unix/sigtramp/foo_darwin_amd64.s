#include "go_asm.h"
#include "textflag.h"

TEXT ·getptr(SB),NOSPLIT,$0
	MOVQ    $runtime·sigtramp(SB), AX
	MOVQ    AX, ret+0(FP)
	RET
