// Code generated from parser/Strictus.g4 by ANTLR 4.7.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 65, 475,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33,
	4, 34, 9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4,
	39, 9, 39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44,
	9, 44, 4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9,
	49, 4, 50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54,
	4, 55, 9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4,
	60, 9, 60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65,
	9, 65, 4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 3, 2, 3,
	2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 8, 3,
	8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12,
	3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3,
	17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 20, 3, 20, 3, 21,
	3, 21, 3, 22, 3, 22, 3, 23, 3, 23, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 3,
	26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 29,
	3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3,
	31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32,
	3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3,
	34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 35, 3, 35,
	3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 36, 3, 36, 3,
	36, 3, 36, 3, 37, 3, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38,
	3, 39, 3, 39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3,
	40, 3, 40, 3, 40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 42,
	3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3,
	43, 3, 43, 3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 3, 44, 3, 45, 3, 45, 3, 45,
	3, 45, 3, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 48, 3,
	48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 49, 3, 49, 3, 49, 3, 49, 3, 49, 3, 50,
	3, 50, 3, 50, 3, 50, 3, 50, 3, 50, 3, 51, 3, 51, 3, 51, 3, 51, 3, 52, 3,
	52, 3, 52, 3, 52, 3, 52, 3, 52, 3, 52, 3, 53, 3, 53, 3, 53, 3, 53, 3, 53,
	3, 54, 3, 54, 7, 54, 355, 10, 54, 12, 54, 14, 54, 358, 11, 54, 3, 55, 5,
	55, 361, 10, 55, 3, 56, 3, 56, 5, 56, 365, 10, 56, 3, 57, 3, 57, 7, 57,
	369, 10, 57, 12, 57, 14, 57, 372, 11, 57, 3, 58, 3, 58, 3, 58, 3, 58, 6,
	58, 378, 10, 58, 13, 58, 14, 58, 379, 3, 59, 3, 59, 3, 59, 3, 59, 6, 59,
	386, 10, 59, 13, 59, 14, 59, 387, 3, 60, 3, 60, 3, 60, 3, 60, 6, 60, 394,
	10, 60, 13, 60, 14, 60, 395, 3, 61, 3, 61, 3, 61, 7, 61, 401, 10, 61, 12,
	61, 14, 61, 404, 11, 61, 3, 62, 3, 62, 7, 62, 408, 10, 62, 12, 62, 14,
	62, 411, 11, 62, 3, 62, 3, 62, 3, 63, 3, 63, 5, 63, 417, 10, 63, 3, 64,
	3, 64, 3, 64, 3, 64, 3, 64, 3, 64, 3, 64, 6, 64, 426, 10, 64, 13, 64, 14,
	64, 427, 3, 64, 3, 64, 5, 64, 432, 10, 64, 3, 65, 3, 65, 3, 66, 6, 66,
	437, 10, 66, 13, 66, 14, 66, 438, 3, 66, 3, 66, 3, 67, 6, 67, 444, 10,
	67, 13, 67, 14, 67, 445, 3, 67, 3, 67, 3, 68, 3, 68, 3, 68, 3, 68, 3, 68,
	7, 68, 455, 10, 68, 12, 68, 14, 68, 458, 11, 68, 3, 68, 3, 68, 3, 68, 3,
	68, 3, 68, 3, 69, 3, 69, 3, 69, 3, 69, 7, 69, 469, 10, 69, 12, 69, 14,
	69, 472, 11, 69, 3, 69, 3, 69, 3, 456, 2, 70, 3, 3, 5, 4, 7, 5, 9, 6, 11,
	7, 13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 16,
	31, 17, 33, 18, 35, 19, 37, 20, 39, 21, 41, 22, 43, 23, 45, 24, 47, 25,
	49, 26, 51, 27, 53, 28, 55, 29, 57, 30, 59, 31, 61, 32, 63, 33, 65, 34,
	67, 35, 69, 36, 71, 37, 73, 38, 75, 39, 77, 40, 79, 41, 81, 42, 83, 43,
	85, 44, 87, 45, 89, 46, 91, 47, 93, 48, 95, 49, 97, 50, 99, 51, 101, 52,
	103, 53, 105, 54, 107, 55, 109, 2, 111, 2, 113, 56, 115, 57, 117, 58, 119,
	59, 121, 60, 123, 61, 125, 2, 127, 2, 129, 2, 131, 62, 133, 63, 135, 64,
	137, 65, 3, 2, 15, 5, 2, 67, 92, 97, 97, 99, 124, 3, 2, 50, 59, 4, 2, 50,
	59, 97, 97, 4, 2, 50, 51, 97, 97, 4, 2, 50, 57, 97, 97, 6, 2, 50, 59, 67,
	72, 97, 97, 99, 104, 4, 2, 67, 92, 99, 124, 6, 2, 50, 59, 67, 92, 97, 97,
	99, 124, 6, 2, 12, 12, 15, 15, 36, 36, 94, 94, 9, 2, 36, 36, 41, 41, 50,
	50, 94, 94, 112, 112, 116, 116, 118, 118, 5, 2, 50, 59, 67, 72, 99, 104,
	6, 2, 2, 2, 11, 11, 13, 14, 34, 34, 4, 2, 12, 12, 15, 15, 2, 485, 2, 3,
	3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11,
	3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2,
	19, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2,
	2, 27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2, 2, 2, 33, 3, 2, 2,
	2, 2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3, 2, 2, 2, 2, 41, 3, 2,
	2, 2, 2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3, 2, 2, 2, 2, 49, 3,
	2, 2, 2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2, 55, 3, 2, 2, 2, 2, 57,
	3, 2, 2, 2, 2, 59, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2, 2, 63, 3, 2, 2, 2, 2,
	65, 3, 2, 2, 2, 2, 67, 3, 2, 2, 2, 2, 69, 3, 2, 2, 2, 2, 71, 3, 2, 2, 2,
	2, 73, 3, 2, 2, 2, 2, 75, 3, 2, 2, 2, 2, 77, 3, 2, 2, 2, 2, 79, 3, 2, 2,
	2, 2, 81, 3, 2, 2, 2, 2, 83, 3, 2, 2, 2, 2, 85, 3, 2, 2, 2, 2, 87, 3, 2,
	2, 2, 2, 89, 3, 2, 2, 2, 2, 91, 3, 2, 2, 2, 2, 93, 3, 2, 2, 2, 2, 95, 3,
	2, 2, 2, 2, 97, 3, 2, 2, 2, 2, 99, 3, 2, 2, 2, 2, 101, 3, 2, 2, 2, 2, 103,
	3, 2, 2, 2, 2, 105, 3, 2, 2, 2, 2, 107, 3, 2, 2, 2, 2, 113, 3, 2, 2, 2,
	2, 115, 3, 2, 2, 2, 2, 117, 3, 2, 2, 2, 2, 119, 3, 2, 2, 2, 2, 121, 3,
	2, 2, 2, 2, 123, 3, 2, 2, 2, 2, 131, 3, 2, 2, 2, 2, 133, 3, 2, 2, 2, 2,
	135, 3, 2, 2, 2, 2, 137, 3, 2, 2, 2, 3, 139, 3, 2, 2, 2, 5, 141, 3, 2,
	2, 2, 7, 143, 3, 2, 2, 2, 9, 145, 3, 2, 2, 2, 11, 147, 3, 2, 2, 2, 13,
	149, 3, 2, 2, 2, 15, 151, 3, 2, 2, 2, 17, 153, 3, 2, 2, 2, 19, 155, 3,
	2, 2, 2, 21, 158, 3, 2, 2, 2, 23, 161, 3, 2, 2, 2, 25, 163, 3, 2, 2, 2,
	27, 166, 3, 2, 2, 2, 29, 169, 3, 2, 2, 2, 31, 171, 3, 2, 2, 2, 33, 173,
	3, 2, 2, 2, 35, 176, 3, 2, 2, 2, 37, 179, 3, 2, 2, 2, 39, 181, 3, 2, 2,
	2, 41, 183, 3, 2, 2, 2, 43, 185, 3, 2, 2, 2, 45, 187, 3, 2, 2, 2, 47, 189,
	3, 2, 2, 2, 49, 191, 3, 2, 2, 2, 51, 194, 3, 2, 2, 2, 53, 196, 3, 2, 2,
	2, 55, 200, 3, 2, 2, 2, 57, 204, 3, 2, 2, 2, 59, 206, 3, 2, 2, 2, 61, 208,
	3, 2, 2, 2, 63, 220, 3, 2, 2, 2, 65, 227, 3, 2, 2, 2, 67, 236, 3, 2, 2,
	2, 69, 245, 3, 2, 2, 2, 71, 255, 3, 2, 2, 2, 73, 259, 3, 2, 2, 2, 75, 263,
	3, 2, 2, 2, 77, 268, 3, 2, 2, 2, 79, 272, 3, 2, 2, 2, 81, 281, 3, 2, 2,
	2, 83, 288, 3, 2, 2, 2, 85, 294, 3, 2, 2, 2, 87, 303, 3, 2, 2, 2, 89, 307,
	3, 2, 2, 2, 91, 311, 3, 2, 2, 2, 93, 314, 3, 2, 2, 2, 95, 319, 3, 2, 2,
	2, 97, 325, 3, 2, 2, 2, 99, 330, 3, 2, 2, 2, 101, 336, 3, 2, 2, 2, 103,
	340, 3, 2, 2, 2, 105, 347, 3, 2, 2, 2, 107, 352, 3, 2, 2, 2, 109, 360,
	3, 2, 2, 2, 111, 364, 3, 2, 2, 2, 113, 366, 3, 2, 2, 2, 115, 373, 3, 2,
	2, 2, 117, 381, 3, 2, 2, 2, 119, 389, 3, 2, 2, 2, 121, 397, 3, 2, 2, 2,
	123, 405, 3, 2, 2, 2, 125, 416, 3, 2, 2, 2, 127, 431, 3, 2, 2, 2, 129,
	433, 3, 2, 2, 2, 131, 436, 3, 2, 2, 2, 133, 443, 3, 2, 2, 2, 135, 449,
	3, 2, 2, 2, 137, 464, 3, 2, 2, 2, 139, 140, 7, 61, 2, 2, 140, 4, 3, 2,
	2, 2, 141, 142, 7, 46, 2, 2, 142, 6, 3, 2, 2, 2, 143, 144, 7, 125, 2, 2,
	144, 8, 3, 2, 2, 2, 145, 146, 7, 127, 2, 2, 146, 10, 3, 2, 2, 2, 147, 148,
	7, 60, 2, 2, 148, 12, 3, 2, 2, 2, 149, 150, 7, 93, 2, 2, 150, 14, 3, 2,
	2, 2, 151, 152, 7, 95, 2, 2, 152, 16, 3, 2, 2, 2, 153, 154, 7, 63, 2, 2,
	154, 18, 3, 2, 2, 2, 155, 156, 7, 126, 2, 2, 156, 157, 7, 126, 2, 2, 157,
	20, 3, 2, 2, 2, 158, 159, 7, 40, 2, 2, 159, 160, 7, 40, 2, 2, 160, 22,
	3, 2, 2, 2, 161, 162, 7, 48, 2, 2, 162, 24, 3, 2, 2, 2, 163, 164, 7, 63,
	2, 2, 164, 165, 7, 63, 2, 2, 165, 26, 3, 2, 2, 2, 166, 167, 7, 35, 2, 2,
	167, 168, 7, 63, 2, 2, 168, 28, 3, 2, 2, 2, 169, 170, 7, 62, 2, 2, 170,
	30, 3, 2, 2, 2, 171, 172, 7, 64, 2, 2, 172, 32, 3, 2, 2, 2, 173, 174, 7,
	62, 2, 2, 174, 175, 7, 63, 2, 2, 175, 34, 3, 2, 2, 2, 176, 177, 7, 64,
	2, 2, 177, 178, 7, 63, 2, 2, 178, 36, 3, 2, 2, 2, 179, 180, 7, 45, 2, 2,
	180, 38, 3, 2, 2, 2, 181, 182, 7, 47, 2, 2, 182, 40, 3, 2, 2, 2, 183, 184,
	7, 44, 2, 2, 184, 42, 3, 2, 2, 2, 185, 186, 7, 49, 2, 2, 186, 44, 3, 2,
	2, 2, 187, 188, 7, 39, 2, 2, 188, 46, 3, 2, 2, 2, 189, 190, 7, 35, 2, 2,
	190, 48, 3, 2, 2, 2, 191, 192, 7, 62, 2, 2, 192, 193, 7, 47, 2, 2, 193,
	50, 3, 2, 2, 2, 194, 195, 7, 65, 2, 2, 195, 52, 3, 2, 2, 2, 196, 197, 5,
	131, 66, 2, 197, 198, 7, 65, 2, 2, 198, 199, 7, 65, 2, 2, 199, 54, 3, 2,
	2, 2, 200, 201, 7, 99, 2, 2, 201, 202, 7, 117, 2, 2, 202, 203, 7, 65, 2,
	2, 203, 56, 3, 2, 2, 2, 204, 205, 7, 42, 2, 2, 205, 58, 3, 2, 2, 2, 206,
	207, 7, 43, 2, 2, 207, 60, 3, 2, 2, 2, 208, 209, 7, 118, 2, 2, 209, 210,
	7, 116, 2, 2, 210, 211, 7, 99, 2, 2, 211, 212, 7, 112, 2, 2, 212, 213,
	7, 117, 2, 2, 213, 214, 7, 99, 2, 2, 214, 215, 7, 101, 2, 2, 215, 216,
	7, 118, 2, 2, 216, 217, 7, 107, 2, 2, 217, 218, 7, 113, 2, 2, 218, 219,
	7, 112, 2, 2, 219, 62, 3, 2, 2, 2, 220, 221, 7, 117, 2, 2, 221, 222, 7,
	118, 2, 2, 222, 223, 7, 116, 2, 2, 223, 224, 7, 119, 2, 2, 224, 225, 7,
	101, 2, 2, 225, 226, 7, 118, 2, 2, 226, 64, 3, 2, 2, 2, 227, 228, 7, 116,
	2, 2, 228, 229, 7, 103, 2, 2, 229, 230, 7, 117, 2, 2, 230, 231, 7, 113,
	2, 2, 231, 232, 7, 119, 2, 2, 232, 233, 7, 116, 2, 2, 233, 234, 7, 101,
	2, 2, 234, 235, 7, 103, 2, 2, 235, 66, 3, 2, 2, 2, 236, 237, 7, 101, 2,
	2, 237, 238, 7, 113, 2, 2, 238, 239, 7, 112, 2, 2, 239, 240, 7, 118, 2,
	2, 240, 241, 7, 116, 2, 2, 241, 242, 7, 99, 2, 2, 242, 243, 7, 101, 2,
	2, 243, 244, 7, 118, 2, 2, 244, 68, 3, 2, 2, 2, 245, 246, 7, 107, 2, 2,
	246, 247, 7, 112, 2, 2, 247, 248, 7, 118, 2, 2, 248, 249, 7, 103, 2, 2,
	249, 250, 7, 116, 2, 2, 250, 251, 7, 104, 2, 2, 251, 252, 7, 99, 2, 2,
	252, 253, 7, 101, 2, 2, 253, 254, 7, 103, 2, 2, 254, 70, 3, 2, 2, 2, 255,
	256, 7, 104, 2, 2, 256, 257, 7, 119, 2, 2, 257, 258, 7, 112, 2, 2, 258,
	72, 3, 2, 2, 2, 259, 260, 7, 114, 2, 2, 260, 261, 7, 116, 2, 2, 261, 262,
	7, 103, 2, 2, 262, 74, 3, 2, 2, 2, 263, 264, 7, 114, 2, 2, 264, 265, 7,
	113, 2, 2, 265, 266, 7, 117, 2, 2, 266, 267, 7, 118, 2, 2, 267, 76, 3,
	2, 2, 2, 268, 269, 7, 114, 2, 2, 269, 270, 7, 119, 2, 2, 270, 271, 7, 100,
	2, 2, 271, 78, 3, 2, 2, 2, 272, 273, 7, 114, 2, 2, 273, 274, 7, 119, 2,
	2, 274, 275, 7, 100, 2, 2, 275, 276, 7, 42, 2, 2, 276, 277, 7, 117, 2,
	2, 277, 278, 7, 103, 2, 2, 278, 279, 7, 118, 2, 2, 279, 280, 7, 43, 2,
	2, 280, 80, 3, 2, 2, 2, 281, 282, 7, 116, 2, 2, 282, 283, 7, 103, 2, 2,
	283, 284, 7, 118, 2, 2, 284, 285, 7, 119, 2, 2, 285, 286, 7, 116, 2, 2,
	286, 287, 7, 112, 2, 2, 287, 82, 3, 2, 2, 2, 288, 289, 7, 100, 2, 2, 289,
	290, 7, 116, 2, 2, 290, 291, 7, 103, 2, 2, 291, 292, 7, 99, 2, 2, 292,
	293, 7, 109, 2, 2, 293, 84, 3, 2, 2, 2, 294, 295, 7, 101, 2, 2, 295, 296,
	7, 113, 2, 2, 296, 297, 7, 112, 2, 2, 297, 298, 7, 118, 2, 2, 298, 299,
	7, 107, 2, 2, 299, 300, 7, 112, 2, 2, 300, 301, 7, 119, 2, 2, 301, 302,
	7, 103, 2, 2, 302, 86, 3, 2, 2, 2, 303, 304, 7, 110, 2, 2, 304, 305, 7,
	103, 2, 2, 305, 306, 7, 118, 2, 2, 306, 88, 3, 2, 2, 2, 307, 308, 7, 120,
	2, 2, 308, 309, 7, 99, 2, 2, 309, 310, 7, 116, 2, 2, 310, 90, 3, 2, 2,
	2, 311, 312, 7, 107, 2, 2, 312, 313, 7, 104, 2, 2, 313, 92, 3, 2, 2, 2,
	314, 315, 7, 103, 2, 2, 315, 316, 7, 110, 2, 2, 316, 317, 7, 117, 2, 2,
	317, 318, 7, 103, 2, 2, 318, 94, 3, 2, 2, 2, 319, 320, 7, 121, 2, 2, 320,
	321, 7, 106, 2, 2, 321, 322, 7, 107, 2, 2, 322, 323, 7, 110, 2, 2, 323,
	324, 7, 103, 2, 2, 324, 96, 3, 2, 2, 2, 325, 326, 7, 118, 2, 2, 326, 327,
	7, 116, 2, 2, 327, 328, 7, 119, 2, 2, 328, 329, 7, 103, 2, 2, 329, 98,
	3, 2, 2, 2, 330, 331, 7, 104, 2, 2, 331, 332, 7, 99, 2, 2, 332, 333, 7,
	110, 2, 2, 333, 334, 7, 117, 2, 2, 334, 335, 7, 103, 2, 2, 335, 100, 3,
	2, 2, 2, 336, 337, 7, 112, 2, 2, 337, 338, 7, 107, 2, 2, 338, 339, 7, 110,
	2, 2, 339, 102, 3, 2, 2, 2, 340, 341, 7, 107, 2, 2, 341, 342, 7, 111, 2,
	2, 342, 343, 7, 114, 2, 2, 343, 344, 7, 113, 2, 2, 344, 345, 7, 116, 2,
	2, 345, 346, 7, 118, 2, 2, 346, 104, 3, 2, 2, 2, 347, 348, 7, 104, 2, 2,
	348, 349, 7, 116, 2, 2, 349, 350, 7, 113, 2, 2, 350, 351, 7, 111, 2, 2,
	351, 106, 3, 2, 2, 2, 352, 356, 5, 109, 55, 2, 353, 355, 5, 111, 56, 2,
	354, 353, 3, 2, 2, 2, 355, 358, 3, 2, 2, 2, 356, 354, 3, 2, 2, 2, 356,
	357, 3, 2, 2, 2, 357, 108, 3, 2, 2, 2, 358, 356, 3, 2, 2, 2, 359, 361,
	9, 2, 2, 2, 360, 359, 3, 2, 2, 2, 361, 110, 3, 2, 2, 2, 362, 365, 9, 3,
	2, 2, 363, 365, 5, 109, 55, 2, 364, 362, 3, 2, 2, 2, 364, 363, 3, 2, 2,
	2, 365, 112, 3, 2, 2, 2, 366, 370, 9, 3, 2, 2, 367, 369, 9, 4, 2, 2, 368,
	367, 3, 2, 2, 2, 369, 372, 3, 2, 2, 2, 370, 368, 3, 2, 2, 2, 370, 371,
	3, 2, 2, 2, 371, 114, 3, 2, 2, 2, 372, 370, 3, 2, 2, 2, 373, 374, 7, 50,
	2, 2, 374, 375, 7, 100, 2, 2, 375, 377, 3, 2, 2, 2, 376, 378, 9, 5, 2,
	2, 377, 376, 3, 2, 2, 2, 378, 379, 3, 2, 2, 2, 379, 377, 3, 2, 2, 2, 379,
	380, 3, 2, 2, 2, 380, 116, 3, 2, 2, 2, 381, 382, 7, 50, 2, 2, 382, 383,
	7, 113, 2, 2, 383, 385, 3, 2, 2, 2, 384, 386, 9, 6, 2, 2, 385, 384, 3,
	2, 2, 2, 386, 387, 3, 2, 2, 2, 387, 385, 3, 2, 2, 2, 387, 388, 3, 2, 2,
	2, 388, 118, 3, 2, 2, 2, 389, 390, 7, 50, 2, 2, 390, 391, 7, 122, 2, 2,
	391, 393, 3, 2, 2, 2, 392, 394, 9, 7, 2, 2, 393, 392, 3, 2, 2, 2, 394,
	395, 3, 2, 2, 2, 395, 393, 3, 2, 2, 2, 395, 396, 3, 2, 2, 2, 396, 120,
	3, 2, 2, 2, 397, 398, 7, 50, 2, 2, 398, 402, 9, 8, 2, 2, 399, 401, 9, 9,
	2, 2, 400, 399, 3, 2, 2, 2, 401, 404, 3, 2, 2, 2, 402, 400, 3, 2, 2, 2,
	402, 403, 3, 2, 2, 2, 403, 122, 3, 2, 2, 2, 404, 402, 3, 2, 2, 2, 405,
	409, 7, 36, 2, 2, 406, 408, 5, 125, 63, 2, 407, 406, 3, 2, 2, 2, 408, 411,
	3, 2, 2, 2, 409, 407, 3, 2, 2, 2, 409, 410, 3, 2, 2, 2, 410, 412, 3, 2,
	2, 2, 411, 409, 3, 2, 2, 2, 412, 413, 7, 36, 2, 2, 413, 124, 3, 2, 2, 2,
	414, 417, 5, 127, 64, 2, 415, 417, 10, 10, 2, 2, 416, 414, 3, 2, 2, 2,
	416, 415, 3, 2, 2, 2, 417, 126, 3, 2, 2, 2, 418, 419, 7, 94, 2, 2, 419,
	432, 9, 11, 2, 2, 420, 421, 7, 94, 2, 2, 421, 422, 7, 119, 2, 2, 422, 423,
	3, 2, 2, 2, 423, 425, 7, 125, 2, 2, 424, 426, 5, 129, 65, 2, 425, 424,
	3, 2, 2, 2, 426, 427, 3, 2, 2, 2, 427, 425, 3, 2, 2, 2, 427, 428, 3, 2,
	2, 2, 428, 429, 3, 2, 2, 2, 429, 430, 7, 127, 2, 2, 430, 432, 3, 2, 2,
	2, 431, 418, 3, 2, 2, 2, 431, 420, 3, 2, 2, 2, 432, 128, 3, 2, 2, 2, 433,
	434, 9, 12, 2, 2, 434, 130, 3, 2, 2, 2, 435, 437, 9, 13, 2, 2, 436, 435,
	3, 2, 2, 2, 437, 438, 3, 2, 2, 2, 438, 436, 3, 2, 2, 2, 438, 439, 3, 2,
	2, 2, 439, 440, 3, 2, 2, 2, 440, 441, 8, 66, 2, 2, 441, 132, 3, 2, 2, 2,
	442, 444, 9, 14, 2, 2, 443, 442, 3, 2, 2, 2, 444, 445, 3, 2, 2, 2, 445,
	443, 3, 2, 2, 2, 445, 446, 3, 2, 2, 2, 446, 447, 3, 2, 2, 2, 447, 448,
	8, 67, 2, 2, 448, 134, 3, 2, 2, 2, 449, 450, 7, 49, 2, 2, 450, 451, 7,
	44, 2, 2, 451, 456, 3, 2, 2, 2, 452, 455, 5, 135, 68, 2, 453, 455, 11,
	2, 2, 2, 454, 452, 3, 2, 2, 2, 454, 453, 3, 2, 2, 2, 455, 458, 3, 2, 2,
	2, 456, 457, 3, 2, 2, 2, 456, 454, 3, 2, 2, 2, 457, 459, 3, 2, 2, 2, 458,
	456, 3, 2, 2, 2, 459, 460, 7, 44, 2, 2, 460, 461, 7, 49, 2, 2, 461, 462,
	3, 2, 2, 2, 462, 463, 8, 68, 2, 2, 463, 136, 3, 2, 2, 2, 464, 465, 7, 49,
	2, 2, 465, 466, 7, 49, 2, 2, 466, 470, 3, 2, 2, 2, 467, 469, 10, 14, 2,
	2, 468, 467, 3, 2, 2, 2, 469, 472, 3, 2, 2, 2, 470, 468, 3, 2, 2, 2, 470,
	471, 3, 2, 2, 2, 471, 473, 3, 2, 2, 2, 472, 470, 3, 2, 2, 2, 473, 474,
	8, 69, 2, 2, 474, 138, 3, 2, 2, 2, 20, 2, 356, 360, 364, 370, 379, 387,
	395, 402, 409, 416, 427, 431, 438, 445, 454, 456, 470, 3, 2, 3, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "';'", "','", "'{'", "'}'", "':'", "'['", "']'", "'='", "'||'", "'&&'",
	"'.'", "'=='", "'!='", "'<'", "'>'", "'<='", "'>='", "'+'", "'-'", "'*'",
	"'/'", "'%'", "'!'", "'<-'", "'?'", "", "'as?'", "'('", "')'", "'transaction'",
	"'struct'", "'resource'", "'contract'", "'interface'", "'fun'", "'pre'",
	"'post'", "'pub'", "'pub(set)'", "'return'", "'break'", "'continue'", "'let'",
	"'var'", "'if'", "'else'", "'while'", "'true'", "'false'", "'nil'", "'import'",
	"'from'",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "Equal", "Unequal", "Less",
	"Greater", "LessEqual", "GreaterEqual", "Plus", "Minus", "Mul", "Div",
	"Mod", "Negate", "Move", "Optional", "NilCoalescing", "FailableDowncasting",
	"OpenParen", "CloseParen", "Transaction", "Struct", "Resource", "Contract",
	"Interface", "Fun", "Pre", "Post", "Pub", "PubSet", "Return", "Break",
	"Continue", "Let", "Var", "If", "Else", "While", "True", "False", "Nil",
	"Import", "From", "Identifier", "DecimalLiteral", "BinaryLiteral", "OctalLiteral",
	"HexadecimalLiteral", "InvalidNumberLiteral", "StringLiteral", "WS", "Terminator",
	"BlockComment", "LineComment",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "T__8",
	"T__9", "T__10", "Equal", "Unequal", "Less", "Greater", "LessEqual", "GreaterEqual",
	"Plus", "Minus", "Mul", "Div", "Mod", "Negate", "Move", "Optional", "NilCoalescing",
	"FailableDowncasting", "OpenParen", "CloseParen", "Transaction", "Struct",
	"Resource", "Contract", "Interface", "Fun", "Pre", "Post", "Pub", "PubSet",
	"Return", "Break", "Continue", "Let", "Var", "If", "Else", "While", "True",
	"False", "Nil", "Import", "From", "Identifier", "IdentifierHead", "IdentifierCharacter",
	"DecimalLiteral", "BinaryLiteral", "OctalLiteral", "HexadecimalLiteral",
	"InvalidNumberLiteral", "StringLiteral", "QuotedText", "EscapedCharacter",
	"HexadecimalDigit", "WS", "Terminator", "BlockComment", "LineComment",
}

type StrictusLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewStrictusLexer(input antlr.CharStream) *StrictusLexer {

	l := new(StrictusLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Strictus.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// StrictusLexer tokens.
const (
	StrictusLexerT__0                 = 1
	StrictusLexerT__1                 = 2
	StrictusLexerT__2                 = 3
	StrictusLexerT__3                 = 4
	StrictusLexerT__4                 = 5
	StrictusLexerT__5                 = 6
	StrictusLexerT__6                 = 7
	StrictusLexerT__7                 = 8
	StrictusLexerT__8                 = 9
	StrictusLexerT__9                 = 10
	StrictusLexerT__10                = 11
	StrictusLexerEqual                = 12
	StrictusLexerUnequal              = 13
	StrictusLexerLess                 = 14
	StrictusLexerGreater              = 15
	StrictusLexerLessEqual            = 16
	StrictusLexerGreaterEqual         = 17
	StrictusLexerPlus                 = 18
	StrictusLexerMinus                = 19
	StrictusLexerMul                  = 20
	StrictusLexerDiv                  = 21
	StrictusLexerMod                  = 22
	StrictusLexerNegate               = 23
	StrictusLexerMove                 = 24
	StrictusLexerOptional             = 25
	StrictusLexerNilCoalescing        = 26
	StrictusLexerFailableDowncasting  = 27
	StrictusLexerOpenParen            = 28
	StrictusLexerCloseParen           = 29
	StrictusLexerTransaction          = 30
	StrictusLexerStruct               = 31
	StrictusLexerResource             = 32
	StrictusLexerContract             = 33
	StrictusLexerInterface            = 34
	StrictusLexerFun                  = 35
	StrictusLexerPre                  = 36
	StrictusLexerPost                 = 37
	StrictusLexerPub                  = 38
	StrictusLexerPubSet               = 39
	StrictusLexerReturn               = 40
	StrictusLexerBreak                = 41
	StrictusLexerContinue             = 42
	StrictusLexerLet                  = 43
	StrictusLexerVar                  = 44
	StrictusLexerIf                   = 45
	StrictusLexerElse                 = 46
	StrictusLexerWhile                = 47
	StrictusLexerTrue                 = 48
	StrictusLexerFalse                = 49
	StrictusLexerNil                  = 50
	StrictusLexerImport               = 51
	StrictusLexerFrom                 = 52
	StrictusLexerIdentifier           = 53
	StrictusLexerDecimalLiteral       = 54
	StrictusLexerBinaryLiteral        = 55
	StrictusLexerOctalLiteral         = 56
	StrictusLexerHexadecimalLiteral   = 57
	StrictusLexerInvalidNumberLiteral = 58
	StrictusLexerStringLiteral        = 59
	StrictusLexerWS                   = 60
	StrictusLexerTerminator           = 61
	StrictusLexerBlockComment         = 62
	StrictusLexerLineComment          = 63
)
