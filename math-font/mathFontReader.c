#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <ft2build.h>
#include <freetype/ftadvanc.h>
#include <freetype/freetype.h>
#include <freetype/ftbbox.h>

// TODO: compress glyphs into a custom ttf file
/*const FT_ULong math_glyphs[] = {
  // basic latin
  0x00028, // parens open
  0x00029, // parens close
  0x0002b, // non-math plus
  0x0002c, // comma separator
	0x0002d, // non-math minus
	0x0002e, // period 
	0x00030, // 0
	0x00031, // 1
	0x00032, // 2
	0x00033, // 3
	0x00034, // 4
	0x00035, // 5
	0x00036, // 6
	0x00037, // 7
	0x00038, // 8
	0x00039, // 9
  0x0003a, // colon separator
  0x0003b, // semicolon separator
	0x00041, // A
	0x00042, // B
	0x00043, // C
	0x00044, // D
	0x00045, // E
	0x00046, // F
	0x00047, // G
	0x00048, // H
	0x00049, // I
	0x0004a, // J
	0x0004b, // K
	0x0004c, // L
	0x0004d, // M
	0x0004e, // N
	0x0004f, // O
	0x00050, // P
	0x00051, // Q
	0x00052, // R
	0x00053, // S
	0x00054, // T
	0x00055, // U
	0x00056, // V
	0x00057, // W
	0x00058, // X
	0x00059, // Y
	0x0005a, // Z
	0x00061, // a
	0x00062, // b
	0x00063, // c
	0x00064, // d
	0x00065, // e
	0x00066, // f
	0x00067, // g
	0x00068, // h
	0x00069, // i
	0x0006a, // j
	0x0006b, // k
	0x0006c, // l
	0x0006d, // m
	0x0006e, // n
	0x0006f, // o
	0x00070, // p
	0x00071, // q
	0x00072, // r
	0x00073, // s
	0x00074, // t
	0x00075, // u
	0x00076, // v
	0x00077, // w
	0x00078, // x
	0x00079, // y
	0x0007a, // z

  // letterlike symbols
  0x0210e, // planck (i.e. h_it)

  // mathematical operators
  0x02202, // partialdiff
  0x02206, // Delta.math
  0x02207, // gradient (i.e. nabla)
  0x0220f, // Product
  0x02211, // sum
  0x02212, // minus
  0x0221e, // infty
  0x02223, // pipe
  0x02225, // double pipe
  0x0222b, // single integral
  0x0222c, // double integral
  0x0222d, // triple integral
  0x02248, // approx

  // halfwidth and fullwidth forms
  0x0ff0b, // (hopefully) wider plus

  // italic latin
  0x1d434, // A_it
  0x1d435, // B_it
  0x1d436, // C_it
  0x1d437, // D_it
  0x1d438, // E_it
  0x1d439, // F_it
  0x1d43a, // G_it
  0x1d43b, // H_it
  0x1d43c, // I_it
  0x1d43d, // J_it
  0x1d43e, // K_it
  0x1d43f, // L_it
  0x1d440, // M_it
  0x1d441, // N_it
  0x1d442, // O_it
  0x1d443, // P_it
  0x1d444, // Q_it
  0x1d445, // R_it
  0x1d446, // S_it
  0x1d447, // T_it
  0x1d448, // U_it
  0x1d449, // V_it
  0x1d44a, // W_it
  0x1d44b, // X_it
  0x1d44c, // Y_it
  0x1d44d, // Z_it
  0x1d44e, // a_it
  0x1d44f, // b_it
  0x1d450, // c_it
  0x1d451, // d_it
  0x1d452, // e_it
  0x1d453, // f_it
  0x1d454, // g_it
  0x1d456, // i_it
  0x1d457, // j_it
  0x1d458, // k_it
  0x1d459, // l_it
  0x1d45a, // m_it
  0x1d45b, // n_it
  0x1d45c, // o_it
  0x1d45d, // p_it
  0x1d45e, // q_it
  0x1d45f, // r_it
  0x1d460, // s_it
  0x1d461, // t_it
  0x1d462, // u_it
  0x1d463, // v_it
  0x1d464, // w_it
  0x1d465, // x_it
  0x1d466, // y_it
  0x1d467, // z_it

  // greek italic
	0x00393, // Gamma
	0x00394, // Delta
	0x00398, // Theta
	0x0039b, // Lambda
	0x0039e, // Xi
	0x003a0, // Pi
	0x003a3, // Sigma
	0x003a5, // Upsilon
	0x003a6, // Phi
	0x003a8, // Psi
	0x003a9, // Omega

	0x1d6fc, // alpha
	0x1d6fd, // beta
	0x1d6fe, // gamma
	0x1d6ff, // delta
	0x1d700, // varepsilon
	0x1d701, // zeta
	0x1d702, // eta
	0x1d703, // theta
	0x1d704, // iota
	0x1d705, // kappa
	0x1d706, // lambda
	0x1d707, // mu
	0x1d708, // nu
	0x1d709, // xi
	0x1d70b, // pi
	0x1d70c, // rho
	0x1d70d, // varsigma
	0x1d70e, // sigma
	0x1d70f, // tau
	0x1d710, // upsilon
	0x1d711, // varphi
	0x1d712, // chi
	0x1d713, // psi
	0x1d714, // omega
	0x1d716, // epsilon
	0x1d718, // varkappa
	0x1d719, // phi
	0x1d71a, // varrho
};*/

// utility functions for working with FT_Fixed, problem is that sometimes they are 16.16, other times 26.6
void print_bits(FT_Fixed fixed) {
  int i;
  for (i = 0; i < 32; i++) {
    printf("%d", (fixed & (1<<i)) ? 1 : 0);
  }
  printf("\n");
}

double fixed1616_to_float(FT_Fixed fixed) {
  int fraction, base;
  double result;
  base = fixed>>16;
  fraction=(fixed) & ((1<<16) - 1);

  result = base + ((double)fraction)/(1<<16);
  return result;
}

double fixed266_to_float(FT_Fixed fixed) {
  int fraction, base;
  double result;
  base = fixed>>6;
  fraction=(fixed) & ((1<<6) - 1);

  result = base + ((double)fraction)/(1<<6);
  return result;
}

void print_advance_widths_table(FT_Face face, FT_ULong *math_glyphs, int n) {
  printf("var AdvanceWidths = map[int]int{\n");

  for (int i = 0 ; i < n; i++) {
    FT_ULong unicode = math_glyphs[i];
    int glyph_index = FT_Get_Char_Index(face, unicode);

    FT_Fixed         advance;
    if (FT_Get_Advance(face, glyph_index, FT_LOAD_NO_SCALE, &advance)) {
      fprintf(stderr, "Error: failed to load glyph advance for %d\n", glyph_index);
      exit(1);
    }

    char glyph_name[256];
    if (FT_Get_Glyph_Name(face, glyph_index, glyph_name, 16)){
      fprintf(stderr, "Error: failed to get glyph name\n");
      exit(1);
    }

    printf("  0x%x :  %d, // %s\n", unicode, advance, glyph_name);
  }

  printf("}\n");
}

void print_bb(FT_Face face, FT_ULong *math_glyphs, int n) {
  printf("var Bounds = map[int]boundingbox.BB{\n");

  for (int i = 0 ; i < n; i++) {
    FT_ULong unicode = math_glyphs[i];
    int glyph_index = FT_Get_Char_Index(face, unicode);

    if (FT_Load_Glyph(face, glyph_index, FT_LOAD_DEFAULT)) {
      fprintf(stderr, "Error: failed to load glyph\n");
      exit(1);
    }

    FT_GlyphSlot slot = face->glyph;

    FT_BBox bb;
    FT_Outline_Get_BBox(&(slot->outline), &bb);

    char glyph_name[256];
    if (FT_Get_Glyph_Name(face, glyph_index, glyph_name, 16)){
      fprintf(stderr, "Error: failed to get glyph name\n");
      exit(1);
    }

    // this bb is different from the one needed in svg, so flip the y
    printf("  0x%x :  boundingbox.NewBB(%g,%g,%g,%g), // %s\n", unicode, 
        fixed266_to_float(bb.xMin), 
        -fixed266_to_float(bb.yMax), 
        fixed266_to_float(bb.xMax), 
        -fixed266_to_float(bb.yMin), glyph_name);
  }

  printf("}\n");
}

int main(int argc, char* argv[]) {
  if (argc < 4) {
    printf("Error: bad number of arguments");
    exit( 1 );
  }

  char* ttf_path = argv[1];
  char* package_name = argv[2];

  // number of glyphs
  int n = argc - 3;
  FT_ULong *math_glyphs = (FT_ULong*)malloc(n*sizeof(FT_ULong));

  // read the glyphs from the command line
  for (int i = 0; i < n; i++) {
    math_glyphs[i] = (FT_ULong)strtol(argv[i+3], NULL, 0);
  }

  FT_Library library;
  if (FT_Init_FreeType(&library)) {
    fprintf(stderr, "Error: failed init library");
    exit(1);
  }

  FT_Face face;
  if (FT_New_Face(library, ttf_path, 0 , &face)) {
    fprintf(stderr, "Error: failed to read face font");
    exit(1);
  }

  printf("package %s\n\n", package_name);
  printf("import ( \"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox\" )\n\n");
  printf("var UnitsPerEm = %d\n\n", face->units_per_EM);

  if (FT_Set_Char_Size(face, 1000 << 6, 1000 << 6, 72, 72)) {
    printf("Error: failed to set char size");
    exit(1);
  }


  print_advance_widths_table(face, math_glyphs, n);

  print_bb(face, math_glyphs, n);

  if (FT_Done_Face(face)) {
    fprintf(stderr, "Error: failed to close face\n");
    exit(1);
  }

  return 0;
}
