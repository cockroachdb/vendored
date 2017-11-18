#include "go_asm.h"
#include "textflag.h"

TEXT ·getptr(SB),NOSPLIT,$0
	MOVD    $runtime·sigtramp(SB), R0
	MOVD    R0, ret+0(FP)
	RET
