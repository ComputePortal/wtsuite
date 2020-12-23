package serif

import ( "github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox" )

var UnitsPerEm = 1000

var AdvanceWidths = map[int]int{
  0x23 :  500, // numbersign
  0x28 :  333, // parenleft
  0x29 :  333, // parenright
  0x2a :  500, // asterisk
  0x2b :  564, // plus
  0x2c :  250, // comma
  0x2d :  333, // hyphen
  0x2e :  250, // period
  0x2f :  296, // slash
  0x30 :  500, // zero
  0x31 :  500, // one
  0x32 :  500, // two
  0x33 :  500, // three
  0x34 :  500, // four
  0x35 :  500, // five
  0x36 :  500, // six
  0x37 :  500, // seven
  0x38 :  500, // eight
  0x39 :  500, // nine
  0x3a :  250, // colon
  0x3b :  250, // semicolon
  0x3c :  564, // less
  0x3e :  564, // greater
  0x41 :  721, // A
  0x42 :  631, // B
  0x43 :  670, // C
  0x44 :  719, // D
  0x45 :  610, // E
  0x46 :  564, // F
  0x47 :  722, // G
  0x48 :  714, // H
  0x49 :  327, // I
  0x4a :  385, // J
  0x4b :  709, // K
  0x4c :  611, // L
  0x4d :  881, // M
  0x4e :  725, // N
  0x4f :  724, // O
  0x50 :  576, // P
  0x51 :  723, // Q
  0x52 :  667, // R
  0x53 :  529, // S
  0x54 :  606, // T
  0x55 :  721, // U
  0x56 :  701, // V
  0x57 :  947, // W
  0x58 :  714, // X
  0x59 :  701, // Y
  0x5a :  613, // Z
  0x61 :  435, // a
  0x62 :  500, // b
  0x63 :  444, // c
  0x64 :  499, // d
  0x65 :  444, // e
  0x66 :  373, // f
  0x67 :  467, // g
  0x68 :  498, // h
  0x69 :  278, // i
  0x6a :  348, // j
  0x6b :  513, // k
  0x6c :  258, // l
  0x6d :  779, // m
  0x6e :  489, // n
  0x6f :  491, // o
  0x70 :  500, // p
  0x71 :  499, // q
  0x72 :  345, // r
  0x73 :  367, // s
  0x74 :  283, // t
  0x75 :  490, // u
  0x76 :  468, // v
  0x77 :  683, // w
  0x78 :  482, // x
  0x79 :  471, // y
  0x7a :  417, // z
  0x7b :  480, // braceleft
  0x7c :  200, // bar
  0x7d :  480, // braceright
  0x210e :  490, // planck
  0x2192 :  900, // arrowright
  0x21d2 :  900, // arrowdblright
  0x2202 :  494, // partialdiff
  0x2206 :  612, // Delta.math
  0x2207 :  612, // gradient
  0x220f :  823, // product
  0x2211 :  713, // summation
  0x2212 :  564, // minus
  0x2219 :  333, // bulletoperator
  0x221e :  853, // infinity
  0x2223 :  200, // divides
  0x2225 :  320, // parallel
  0x222b :  456, // integral
  0x222c :  812, // integraldbl
  0x222d :  1153, // integraltrpl
  0x2248 :  636, // approxequal
  0x2260 :  564, // notequal
  0x2264 :  636, // lessequal
  0x2265 :  636, // greaterequal
  0x226a :  900, // uni226A
  0x226b :  899, // uni226B
  0x3bc :  536, // mu
  0xff0b :  600, // .notdef
  0x1d434 :  611, // A_it
  0x1d435 :  611, // B_it
  0x1d436 :  667, // C_it
  0x1d437 :  722, // D_it
  0x1d438 :  604, // E_it
  0x1d439 :  611, // F_it
  0x1d43a :  722, // G_it
  0x1d43b :  722, // H_it
  0x1d43c :  339, // I_it
  0x1d43d :  444, // J_it
  0x1d43e :  652, // K_it
  0x1d43f :  556, // L_it
  0x1d440 :  828, // M_it
  0x1d441 :  657, // N_it
  0x1d442 :  722, // O_it
  0x1d443 :  603, // P_it
  0x1d444 :  722, // Q_it
  0x1d445 :  616, // R_it
  0x1d446 :  500, // S_it
  0x1d447 :  556, // T_it
  0x1d448 :  722, // U_it
  0x1d449 :  611, // V_it
  0x1d44a :  833, // W_it
  0x1d44b :  611, // X_it
  0x1d44c :  556, // Y_it
  0x1d44d :  556, // Z_it
  0x1d44e :  500, // a_it
  0x1d44f :  500, // b_it
  0x1d450 :  444, // c_it
  0x1d451 :  500, // d_it
  0x1d452 :  444, // e_it
  0x1d453 :  278, // f_it
  0x1d454 :  500, // g_it
  0x1d456 :  278, // i_it
  0x1d457 :  278, // j_it
  0x1d458 :  444, // k_it
  0x1d459 :  278, // l_it
  0x1d45a :  722, // m_it
  0x1d45b :  500, // n_it
  0x1d45c :  500, // o_it
  0x1d45d :  500, // p_it
  0x1d45e :  500, // q_it
  0x1d45f :  389, // r_it
  0x1d460 :  389, // s_it
  0x1d461 :  278, // t_it
  0x1d462 :  500, // u_it
  0x1d463 :  444, // v_it
  0x1d464 :  667, // w_it
  0x1d465 :  444, // x_it
  0x1d466 :  444, // y_it
  0x1d467 :  389, // z_it
  0x393 :  569, // Gamma
  0x394 :  660, // Delta
  0x398 :  754, // Theta
  0x39b :  721, // Lambda
  0x39e :  590, // Xi
  0x3a0 :  713, // Pi
  0x3a3 :  603, // Sigma
  0x3a5 :  666, // Upsilon
  0x3a6 :  760, // Phi
  0x3a8 :  788, // Psi
  0x3a9 :  723, // Omega
  0x1d6fc :  564, // alpha_it
  0x1d6fd :  509, // beta_it
  0x1d6fe :  496, // gamma_it
  0x1d6ff :  520, // delta_it
  0x1d700 :  416, // epsilon_it
  0x1d701 :  398, // zeta_it
  0x1d702 :  506, // eta_it
  0x1d703 :  533, // theta_it
  0x1d704 :  270, // iota_it
  0x1d705 :  491, // kappa_it
  0x1d706 :  488, // lamda_it
  0x1d707 :  501, // mu_it
  0x1d708 :  486, // nu_it
  0x1d709 :  430, // xi_it
  0x1d70b :  608, // pi_it
  0x1d70c :  506, // rho_it
  0x1d70d :  423, // finalsigma_it
  0x1d70e :  524, // sigma_it
  0x1d70f :  425, // tau_it
  0x1d710 :  504, // upsilon_it
  0x1d711 :  618, // phi_it
  0x1d712 :  459, // chi_it
  0x1d713 :  693, // psi_it
  0x1d714 :  693, // omega_it
  0x1d716 :  280, // epsilonsymbol_i
  0x1d718 :  534, // kappasymbol_it
  0x1d719 :  640, // phisymbol_it
  0x1d71a :  534, // rhosymbol_it
}
var Bounds = map[int]boundingbox.BB{
  0x23 :  boundingbox.NewBB(5,-662,496,-0), // numbersign
  0x28 :  boundingbox.NewBB(48,-676,304,177), // parenleft
  0x29 :  boundingbox.NewBB(29,-676,285,177), // parenright
  0x2a :  boundingbox.NewBB(69,-676,432,-265), // asterisk
  0x2b :  boundingbox.NewBB(30,-506,534,-0), // plus
  0x2c :  boundingbox.NewBB(56,-102,195,141), // comma
  0x2d :  boundingbox.NewBB(39,-257,285,-194), // hyphen
  0x2e :  boundingbox.NewBB(70,-100,181,11), // period
  0x2f :  boundingbox.NewBB(0,-676,296,14), // slash
  0x30 :  boundingbox.NewBB(24,-676,476,14), // zero
  0x31 :  boundingbox.NewBB(111,-676,394,-0), // one
  0x32 :  boundingbox.NewBB(30,-676,475,-0), // two
  0x33 :  boundingbox.NewBB(43,-676,432,14), // three
  0x34 :  boundingbox.NewBB(12,-676,472,-0), // four
  0x35 :  boundingbox.NewBB(32,-688,438,14), // five
  0x36 :  boundingbox.NewBB(34,-684,468,14), // six
  0x37 :  boundingbox.NewBB(20,-662,449,8), // seven
  0x38 :  boundingbox.NewBB(56,-676,445,14), // eight
  0x39 :  boundingbox.NewBB(30,-676,459,22), // nine
  0x3a :  boundingbox.NewBB(81,-459,192,11), // colon
  0x3b :  boundingbox.NewBB(80,-459,219,141), // semicolon
  0x3c :  boundingbox.NewBB(28,-516,536,10), // less
  0x3e :  boundingbox.NewBB(28,-516,536,10), // greater
  0x41 :  boundingbox.NewBB(15,-674,706,-0), // A
  0x42 :  boundingbox.NewBB(15,-662,591,-0), // B
  0x43 :  boundingbox.NewBB(35,-676,640,14), // C
  0x44 :  boundingbox.NewBB(15,-662,684,-0), // D
  0x45 :  boundingbox.NewBB(15,-662,600,-0), // E
  0x46 :  boundingbox.NewBB(15,-662,549,-0), // F
  0x47 :  boundingbox.NewBB(35,-676,712,14), // G
  0x48 :  boundingbox.NewBB(15,-662,698,-0), // H
  0x49 :  boundingbox.NewBB(15,-662,312,-0), // I
  0x4a :  boundingbox.NewBB(10,-662,370,14), // J
  0x4b :  boundingbox.NewBB(15,-662,704,-0), // K
  0x4c :  boundingbox.NewBB(15,-662,601,-0), // L
  0x4d :  boundingbox.NewBB(15,-662,866,-0), // M
  0x4e :  boundingbox.NewBB(15,-662,710,11), // N
  0x4f :  boundingbox.NewBB(35,-676,689,14), // O
  0x50 :  boundingbox.NewBB(15,-662,541,-0), // P
  0x51 :  boundingbox.NewBB(35,-676,702,178), // Q
  0x52 :  boundingbox.NewBB(15,-662,657,-0), // R
  0x53 :  boundingbox.NewBB(30,-676,479,14), // S
  0x54 :  boundingbox.NewBB(15,-662,591,-0), // T
  0x55 :  boundingbox.NewBB(15,-662,706,14), // U
  0x56 :  boundingbox.NewBB(10,-662,691,11), // V
  0x57 :  boundingbox.NewBB(10,-662,937,11), // W
  0x58 :  boundingbox.NewBB(10,-662,704,-0), // X
  0x59 :  boundingbox.NewBB(10,-662,691,-0), // Y
  0x5a :  boundingbox.NewBB(10,-662,598,-0), // Z
  0x61 :  boundingbox.NewBB(25,-460,430,10), // a
  0x62 :  boundingbox.NewBB(10,-683,475,10), // b
  0x63 :  boundingbox.NewBB(25,-460,412,10), // c
  0x64 :  boundingbox.NewBB(25,-683,489,10), // d
  0x65 :  boundingbox.NewBB(25,-460,424,10), // e
  0x66 :  boundingbox.NewBB(10,-683,373,-0), // f
  0x67 :  boundingbox.NewBB(15,-460,457,218), // g
  0x68 :  boundingbox.NewBB(10,-683,488,-0), // h
  0x69 :  boundingbox.NewBB(16,-683,253,-0), // i
  0x6a :  boundingbox.NewBB(0,-683,264,218), // j
  0x6b :  boundingbox.NewBB(10,-683,508,-0), // k
  0x6c :  boundingbox.NewBB(10,-683,248,-0), // l
  0x6d :  boundingbox.NewBB(10,-460,769,-0), // m
  0x6e :  boundingbox.NewBB(10,-460,479,-0), // n
  0x6f :  boundingbox.NewBB(25,-460,466,10), // o
  0x70 :  boundingbox.NewBB(10,-460,475,217), // p
  0x71 :  boundingbox.NewBB(25,-461,489,217), // q
  0x72 :  boundingbox.NewBB(10,-460,340,-0), // r
  0x73 :  boundingbox.NewBB(35,-459,332,10), // s
  0x74 :  boundingbox.NewBB(17,-579,283,10), // t
  0x75 :  boundingbox.NewBB(10,-450,480,10), // u
  0x76 :  boundingbox.NewBB(5,-450,463,14), // v
  0x77 :  boundingbox.NewBB(5,-450,678,14), // w
  0x78 :  boundingbox.NewBB(10,-450,472,-0), // x
  0x79 :  boundingbox.NewBB(5,-450,466,218), // y
  0x7a :  boundingbox.NewBB(10,-450,401,-0), // z
  0x7b :  boundingbox.NewBB(100,-680,350,181), // braceleft
  0x7c :  boundingbox.NewBB(67,-676,133,14), // bar
  0x7d :  boundingbox.NewBB(130,-680,380,181), // braceright
  0x210e :  boundingbox.NewBB(9,-683,468,9), // planck
  0x2192 :  boundingbox.NewBB(30,-462,870,-52), // arrowright
  0x21d2 :  boundingbox.NewBB(30,-508,870,-3), // arrowdblright
  0x2202 :  boundingbox.NewBB(26,-675,462,10), // partialdiff
  0x2206 :  boundingbox.NewBB(6,-688,608,-0), // Delta.math
  0x2207 :  boundingbox.NewBB(6,-662,608,26), // gradient
  0x220f :  boundingbox.NewBB(25,-751,803,101), // product
  0x2211 :  boundingbox.NewBB(14,-752,695,123), // summation
  0x2212 :  boundingbox.NewBB(30,-286,534,-220), // minus
  0x2219 :  boundingbox.NewBB(66,-354,268,-152), // bulletoperator
  0x221e :  boundingbox.NewBB(52,-440,801,-62), // infinity
  0x2223 :  boundingbox.NewBB(72,-662,128,157), // divides
  0x2225 :  boundingbox.NewBB(72,-662,248,157), // parallel
  0x222b :  boundingbox.NewBB(0,-900,499,200), // integral
  0x222c :  boundingbox.NewBB(0,-900,849,200), // integraldbl
  0x222d :  boundingbox.NewBB(0,-900,1193,200), // integraltrpl
  0x2248 :  boundingbox.NewBB(56,-440,581,-74), // approxequal
  0x2260 :  boundingbox.NewBB(30,-566,534,52), // notequal
  0x2264 :  boundingbox.NewBB(64,-633,574,-1), // lessequal
  0x2265 :  boundingbox.NewBB(62,-641,570,-1), // greaterequal
  0x226a :  boundingbox.NewBB(28,-516,867,10), // uni226A
  0x226b :  boundingbox.NewBB(28,-516,866,10), // uni226B
  0x3bc :  boundingbox.NewBB(65,-451,521,224), // mu
  0xff0b :  boundingbox.NewBB(34,-750,566,71), // .notdef
  0x1d434 :  boundingbox.NewBB(-51,-668,564,-0), // A_it
  0x1d435 :  boundingbox.NewBB(-8,-653,588,-0), // B_it
  0x1d436 :  boundingbox.NewBB(66,-666,689,18), // C_it
  0x1d437 :  boundingbox.NewBB(-8,-653,700,-0), // D_it
  0x1d438 :  boundingbox.NewBB(-8,-653,627,-0), // E_it
  0x1d439 :  boundingbox.NewBB(-8,-653,629,-0), // F_it
  0x1d43a :  boundingbox.NewBB(52,-666,722,18), // G_it
  0x1d43b :  boundingbox.NewBB(-8,-653,767,-0), // H_it
  0x1d43c :  boundingbox.NewBB(-8,-653,384,-0), // I_it
  0x1d43d :  boundingbox.NewBB(-6,-653,491,18), // J_it
  0x1d43e :  boundingbox.NewBB(-8,-653,707,-0), // K_it
  0x1d43f :  boundingbox.NewBB(-8,-653,559,-0), // L_it
  0x1d440 :  boundingbox.NewBB(-18,-653,873,-0), // M_it
  0x1d441 :  boundingbox.NewBB(-18,-653,729,15), // N_it
  0x1d442 :  boundingbox.NewBB(60,-666,699,18), // O_it
  0x1d443 :  boundingbox.NewBB(-8,-653,597,-0), // P_it
  0x1d444 :  boundingbox.NewBB(59,-666,699,183), // Q_it
  0x1d445 :  boundingbox.NewBB(-8,-653,593,-0), // R_it
  0x1d446 :  boundingbox.NewBB(17,-667,508,18), // S_it
  0x1d447 :  boundingbox.NewBB(59,-653,633,-0), // T_it
  0x1d448 :  boundingbox.NewBB(102,-653,765,18), // U_it
  0x1d449 :  boundingbox.NewBB(76,-653,688,18), // V_it
  0x1d44a :  boundingbox.NewBB(71,-653,906,18), // W_it
  0x1d44b :  boundingbox.NewBB(-29,-653,655,-0), // X_it
  0x1d44c :  boundingbox.NewBB(78,-653,633,-0), // Y_it
  0x1d44d :  boundingbox.NewBB(-6,-653,606,-0), // Z_it
  0x1d44e :  boundingbox.NewBB(17,-441,476,11), // a_it
  0x1d44f :  boundingbox.NewBB(23,-683,473,11), // b_it
  0x1d450 :  boundingbox.NewBB(30,-441,425,11), // c_it
  0x1d451 :  boundingbox.NewBB(15,-683,527,13), // d_it
  0x1d452 :  boundingbox.NewBB(31,-441,412,11), // e_it
  0x1d453 :  boundingbox.NewBB(-147,-678,424,207), // f_it
  0x1d454 :  boundingbox.NewBB(8,-441,472,206), // g_it
  0x1d456 :  boundingbox.NewBB(49,-654,264,11), // i_it
  0x1d457 :  boundingbox.NewBB(-124,-654,276,207), // j_it
  0x1d458 :  boundingbox.NewBB(14,-683,461,11), // k_it
  0x1d459 :  boundingbox.NewBB(40,-683,279,11), // l_it
  0x1d45a :  boundingbox.NewBB(12,-441,704,9), // m_it
  0x1d45b :  boundingbox.NewBB(14,-441,474,9), // n_it
  0x1d45c :  boundingbox.NewBB(27,-441,468,11), // o_it
  0x1d45d :  boundingbox.NewBB(-75,-441,469,205), // p_it
  0x1d45e :  boundingbox.NewBB(25,-441,483,209), // q_it
  0x1d45f :  boundingbox.NewBB(45,-441,412,-0), // r_it
  0x1d460 :  boundingbox.NewBB(16,-442,366,13), // s_it
  0x1d461 :  boundingbox.NewBB(37,-546,296,11), // t_it
  0x1d462 :  boundingbox.NewBB(42,-441,475,11), // u_it
  0x1d463 :  boundingbox.NewBB(21,-441,426,18), // v_it
  0x1d464 :  boundingbox.NewBB(16,-441,648,18), // w_it
  0x1d465 :  boundingbox.NewBB(-27,-441,447,11), // x_it
  0x1d466 :  boundingbox.NewBB(-24,-441,426,206), // y_it
  0x1d467 :  boundingbox.NewBB(-2,-428,380,81), // z_it
  0x393 :  boundingbox.NewBB(15,-662,549,-0), // Gamma
  0x394 :  boundingbox.NewBB(10,-674,640,-0), // Delta
  0x398 :  boundingbox.NewBB(35,-673,719,12), // Theta
  0x39b :  boundingbox.NewBB(15,-674,706,-0), // Lambda
  0x39e :  boundingbox.NewBB(40,-662,550,-0), // Xi
  0x3a0 :  boundingbox.NewBB(15,-662,698,-0), // Pi
  0x3a3 :  boundingbox.NewBB(15,-661,563,-0), // Sigma
  0x3a5 :  boundingbox.NewBB(15,-672,651,-0), // Upsilon
  0x3a6 :  boundingbox.NewBB(35,-662,725,-0), // Phi
  0x3a8 :  boundingbox.NewBB(15,-668,773,-0), // Psi
  0x3a9 :  boundingbox.NewBB(15,-676,708,-0), // Omega
  0x1d6fc :  boundingbox.NewBB(9,-441,527,11), // alpha_it
  0x1d6fd :  boundingbox.NewBB(-34,-674,491,212), // beta_it
  0x1d6fe :  boundingbox.NewBB(34,-441,472,212), // gamma_it
  0x1d6ff :  boundingbox.NewBB(28,-677,488,11), // delta_it
  0x1d700 :  boundingbox.NewBB(-4,-441,386,10), // epsilon_it
  0x1d701 :  boundingbox.NewBB(12,-684,432,178), // zeta_it
  0x1d702 :  boundingbox.NewBB(28,-441,455,212), // eta_it
  0x1d703 :  boundingbox.NewBB(35,-672,507,11), // theta_it
  0x1d704 :  boundingbox.NewBB(3,-432,229,11), // iota_it
  0x1d705 :  boundingbox.NewBB(14,-441,466,10), // kappa_it
  0x1d706 :  boundingbox.NewBB(-11,-678,447,11), // lamda_it
  0x1d707 :  boundingbox.NewBB(-51,-431,452,212), // mu_it
  0x1d708 :  boundingbox.NewBB(12,-441,461,-0), // nu_it
  0x1d709 :  boundingbox.NewBB(0,-684,435,178), // xi_it
  0x1d70b :  boundingbox.NewBB(-22,-451,597,13), // pi_it
  0x1d70c :  boundingbox.NewBB(-12,-441,451,218), // rho_it
  0x1d70d :  boundingbox.NewBB(-8,-441,411,178), // finalsigma_it
  0x1d70e :  boundingbox.NewBB(10,-450,537,11), // sigma_it
  0x1d70f :  boundingbox.NewBB(11,-451,423,13), // tau_it
  0x1d710 :  boundingbox.NewBB(25,-441,459,11), // upsilon_it
  0x1d711 :  boundingbox.NewBB(13,-441,581,205), // phi_it
  0x1d712 :  boundingbox.NewBB(-93,-441,457,205), // chi_it
  0x1d713 :  boundingbox.NewBB(23,-558,654,224), // psi_it
  0x1d714 :  boundingbox.NewBB(11,-461,640,10), // omega_it
  0x1d716 :  boundingbox.NewBB(9,-461,307,10), // epsilonsymbol_i
  0x1d718 :  boundingbox.NewBB(5,-461,481,10), // kappasymbol_it
  0x1d719 :  boundingbox.NewBB(12,-662,582,224), // phisymbol_it
  0x1d71a :  boundingbox.NewBB(-12,-461,470,224), // rhosymbol_it
}
