// generated by stringer -type=TokenType token.go; DO NOT EDIT

package ast

import "fmt"

const _TokenType_name = "T_SPACET_COMMENT_LINET_COMMENT_BLOCKT_SEMICOLONT_COMMAT_IDENTT_URLT_MEDIAT_PAGET_TRUET_FALSET_NULLT_ONLYT_MS_PARAM_NAMET_FUNCTION_NAMET_ID_SELECTORT_CLASS_SELECTORT_TYPE_SELECTORT_UNIVERSAL_SELECTORT_PARENT_SELECTORT_PSEUDO_SELECTORT_FUNCTIONAL_PSEUDOT_INTERPOLATION_SELECTORT_LITERAL_CONCATT_CONCATT_MS_PROGIDT_AND_SELECTORT_DESCENDANT_COMBINATORT_CHILD_COMBINATORT_ADJACENT_SIBLING_COMBINATORT_GENERAL_SIBLING_COMBINATORT_UNICODE_RANGET_IFT_ELSET_ELSE_IFT_INCLUDET_EACHT_WHENT_MIXINT_FUNCTIONT_FORT_FOR_FROMT_FOR_THROUGHT_FOR_TOT_FOR_INT_WHILET_RETURNT_RANGET_CONTENTT_GLOBALT_DEFAULTT_IMPORTANTT_OPTIONALT_FONT_FACET_LOGICAL_NOTT_LOGICAL_ORT_LOGICAL_ANDT_LOGICAL_XORT_NOPT_PLUST_DIVT_MULT_MINUST_MODT_BRACE_OPENT_BRACE_CLOSET_LANG_CODET_BRACKET_OPENT_ATTRIBUTE_NAMET_BRACKET_CLOSET_EQUALT_UNEQUALT_GTT_LTT_GET_LET_ASSIGNT_ATTR_EQUALT_INCLUDE_MATCHT_PREFIX_MATCHT_DASH_MATCHT_SUFFIX_MATCHT_SUBSTRING_MATCHT_VARIABLET_VARIABLE_LENGTH_ARGUMENTST_IMPORTT_AT_RULET_CHARSETT_QQ_STRINGT_Q_STRINGT_UNQUOTE_STRINGT_PAREN_OPENT_PAREN_CLOSET_CONSTANTT_INTEGERT_FLOATT_CDOT_CDCT_UNIT_NONET_UNIT_PERCENTT_UNIT_SECONDT_UNIT_MILLISECONDT_UNIT_EMT_UNIT_EXT_UNIT_CHT_UNIT_REMT_UNIT_CMT_UNIT_INT_UNIT_MMT_UNIT_PCT_UNIT_PTT_UNIT_PXT_UNIT_VHT_UNIT_VWT_UNIT_VMINT_UNIT_VMAXT_UNIT_HZT_UNIT_KHZT_UNIT_DPIT_UNIT_DPCMT_UNIT_DPPXT_UNIT_DEGT_UNIT_GRADT_UNIT_RADT_UNIT_TURNT_PROPERTY_NAME_TOKENT_PROPERTY_VALUET_HEX_COLORT_COLONT_INTERPOLATION_STARTT_INTERPOLATION_INNERT_INTERPOLATION_END"

var _TokenType_index = [...]uint16{0, 7, 21, 36, 47, 54, 61, 66, 73, 79, 85, 92, 98, 104, 119, 134, 147, 163, 178, 198, 215, 232, 251, 275, 291, 299, 310, 324, 347, 365, 394, 422, 437, 441, 447, 456, 465, 471, 477, 484, 494, 499, 509, 522, 530, 538, 545, 553, 560, 569, 577, 586, 597, 607, 618, 631, 643, 656, 669, 674, 680, 685, 690, 697, 702, 714, 727, 738, 752, 768, 783, 790, 799, 803, 807, 811, 815, 823, 835, 850, 864, 876, 890, 907, 917, 944, 952, 961, 970, 981, 991, 1007, 1019, 1032, 1042, 1051, 1058, 1063, 1068, 1079, 1093, 1106, 1124, 1133, 1142, 1151, 1161, 1170, 1179, 1188, 1197, 1206, 1215, 1224, 1233, 1244, 1255, 1264, 1274, 1284, 1295, 1306, 1316, 1327, 1337, 1348, 1369, 1385, 1396, 1403, 1424, 1445, 1464}

func (i TokenType) String() string {
	if i < 0 || i+1 >= TokenType(len(_TokenType_index)) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
