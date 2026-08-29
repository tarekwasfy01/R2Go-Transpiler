package rgo

// do_lapack translates batch_142.c (src/main/lapack.c); R primitive(s): La_qr_cmplx, La_rg, La_rg_cmplx, La_rs, La_rs_cmplx, La_solve_cmplx, La_svd, La_svd_cmplx, La_zgecon, La_ztrcon, La_ztrcon3, qr_coef_cmplx, qr_coef_real, qr_qy_cmplx, qr_qy_real.
func do_lapack(call, op, args, env Value) Value {
	return unsupported("do_lapack", "GNU R LAPACK bridge is not present; Pure-Go linear algebra backend required")
}
