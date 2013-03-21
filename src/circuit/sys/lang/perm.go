package lang

/*

					Ref			PermRef			exportRewrite		importRewrite
	-----------------------------------------------------------------------------
	T				*ref		*permref		=					=
	xptr			=			panic			ptrPtrMsg			n/a
	*ref			=			panic			ptrMsg				n/a
	xpermptr		=			=				permPtrPtrMsg		n/a
	*permref		panic		=				permPtrMsg			n/a
	-----------------------------------------------------------------------------
	*ptrMsg			n/a			n/a				n/a					*ptr
	*ptrPtrMsg		n/a			n/a				n/a					*ptr
	*permPtrMsg		n/a			n/a				n/a					*permptr
	*permPtrPtrMsg	n/a			n/a				n/a					*permptr


	USER VS RUNTIME TYPES

	X     ≈ xptr,		*_ref,		*_ptr,		*ptrMsg,		*ptrPtrMsg
	XPerm ≈ xpermptr,	*_permref,	*_permptr,	*permPtrMsg,	*permPtrPtrMsg

*/
