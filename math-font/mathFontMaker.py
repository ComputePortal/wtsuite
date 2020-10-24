#!/usr/bin/env python3
import sys
import fontTools.subset as ftss
import io
import base64
import subprocess

glyphs = [
# basic latin
0x00023, # hash
0x00028, # parens open
0x00029, # parens close
0x0002a, # non-math multiplication
0x0002b, # non-math plus
0x0002c, # comma separator
0x0002d, # non-math minus
0x0002e, # period 
0x0002f, # forward slash (for division) 
0x00030, # 0
0x00031, # 1
0x00032, # 2
0x00033, # 3
0x00034, # 4
0x00035, # 5
0x00036, # 6
0x00037, # 7
0x00038, # 8
0x00039, # 9
0x0003a, # colon separator
0x0003b, # semicolon separator
0x0003c, # less than
0x0003e, # greater than
0x00041, # A
0x00042, # B
0x00043, # C
0x00044, # D
0x00045, # E
0x00046, # F
0x00047, # G
0x00048, # H
0x00049, # I
0x0004a, # J
0x0004b, # K
0x0004c, # L
0x0004d, # M
0x0004e, # N
0x0004f, # O
0x00050, # P
0x00051, # Q
0x00052, # R
0x00053, # S
0x00054, # T
0x00055, # U
0x00056, # V
0x00057, # W
0x00058, # X
0x00059, # Y
0x0005a, # Z
0x00061, # a
0x00062, # b
0x00063, # c
0x00064, # d
0x00065, # e
0x00066, # f
0x00067, # g
0x00068, # h
0x00069, # i
0x0006a, # j
0x0006b, # k
0x0006c, # l
0x0006d, # m
0x0006e, # n
0x0006f, # o
0x00070, # p
0x00071, # q
0x00072, # r
0x00073, # s
0x00074, # t
0x00075, # u
0x00076, # v
0x00077, # w
0x00078, # x
0x00079, # y
0x0007a, # z
0x0007b, # {
0x0007c, # |
0x0007d, # }

# letterlike symbols
0x0210e, # planck (i.e. h_it)

# mathematical operators
0x02192, # right arrow 1
0x021d2, # right arrow 2
0x02202, # partialdiff
0x02206, # Delta.math
0x02207, # gradient (i.e. nabla)
0x0220f, # Product
0x02211, # sum
0x02212, # minus
0x02219, # dot
0x0221e, # infty
0x02223, # pipe
0x02225, # double pipe
0x0222b, # single integral
0x0222c, # double integral
0x0222d, # triple integral
0x02248, # approx
0x02260, # not equals
0x02264, # less-than or equals
0x02265, # greater than or equals
0x0226a, # much less than
0x0226b, # much greater than

# basic greek (eg. non-italic mu for units)
0x03bc,

# halfwidth and fullwidth forms
0x0ff0b, # (hopefully) wider plus

# italic latin
0x1d434, # A_it
0x1d435, # B_it
0x1d436, # C_it
0x1d437, # D_it
0x1d438, # E_it
0x1d439, # F_it
0x1d43a, # G_it
0x1d43b, # H_it
0x1d43c, # I_it
0x1d43d, # J_it
0x1d43e, # K_it
0x1d43f, # L_it
0x1d440, # M_it
0x1d441, # N_it
0x1d442, # O_it
0x1d443, # P_it
0x1d444, # Q_it
0x1d445, # R_it
0x1d446, # S_it
0x1d447, # T_it
0x1d448, # U_it
0x1d449, # V_it
0x1d44a, # W_it
0x1d44b, # X_it
0x1d44c, # Y_it
0x1d44d, # Z_it
0x1d44e, # a_it
0x1d44f, # b_it
0x1d450, # c_it
0x1d451, # d_it
0x1d452, # e_it
0x1d453, # f_it
0x1d454, # g_it
0x1d456, # i_it
0x1d457, # j_it
0x1d458, # k_it
0x1d459, # l_it
0x1d45a, # m_it
0x1d45b, # n_it
0x1d45c, # o_it
0x1d45d, # p_it
0x1d45e, # q_it
0x1d45f, # r_it
0x1d460, # s_it
0x1d461, # t_it
0x1d462, # u_it
0x1d463, # v_it
0x1d464, # w_it
0x1d465, # x_it
0x1d466, # y_it
0x1d467, # z_it

# greek italic
0x00393, # Gamma
0x00394, # Delta
0x00398, # Theta
0x0039b, # Lambda
0x0039e, # Xi
0x003a0, # Pi
0x003a3, # Sigma
0x003a5, # Upsilon
0x003a6, # Phi
0x003a8, # Psi
0x003a9, # Omega

0x1d6fc, # alpha
0x1d6fd, # beta
0x1d6fe, # gamma
0x1d6ff, # delta
0x1d700, # varepsilon
0x1d701, # zeta
0x1d702, # eta
0x1d703, # theta
0x1d704, # iota
0x1d705, # kappa
0x1d706, # lambda
0x1d707, # mu
0x1d708, # nu
0x1d709, # xi
0x1d70b, # pi
0x1d70c, # rho
0x1d70d, # varsigma
0x1d70e, # sigma
0x1d70f, # tau
0x1d710, # upsilon
0x1d711, # varphi
0x1d712, # chi
0x1d713, # psi
0x1d714, # omega
0x1d716, # epsilon
0x1d718, # varkappa
0x1d719, # phi
0x1d71a, # varrho
]

if len(sys.argv) != 5:
    print("Usage: " + sys.argv[0] + " font-file font-reader dimension-output woff2-output")
    exit()

fname = sys.argv[1]
freader = sys.argv[2]
doutput = sys.argv[3]
woutput = sys.argv[4]

options = ftss.Options()
options.flavor = 'woff2'

subsetter = ftss.Subsetter(options=options)

font = ftss.load_font(fname, options)

subsetter.populate(unicodes=glyphs)

subsetter.subset(font)

out = io.BytesIO()

ftss.save_font(font, out, options)

font.close()

# save the woff2 file
fw = open(woutput, "w")
fw.write("package serif\n")
fw.write("var Woff2Blob = \"")
fw.write(base64.b64encode(out.getvalue()).decode("ascii"))
fw.write("\"")
fw.close()

# save the dimensions file
fd = open(doutput, "w", 1)
subargs = [freader, fname, "serif"]
subargs.extend(map(lambda x: str(x), glyphs))
proc = subprocess.Popen(subargs, stdout=fd)
proc.wait()
fd.close()
