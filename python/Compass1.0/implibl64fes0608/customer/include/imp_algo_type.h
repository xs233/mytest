/* File:	imp_algo_type.h
* Version: 2.0
*
* Author: IMPower Algorithm Team
*
* Description:	Fundamental type definitions for IMPower algorithms
*
* Copyright (C) 2009-2016, IMPower Technologies, all rights reserved.
*
* Change log:
*
* 20160830 Clear redundant definitions
*
*
*/
#ifndef _IMP_ALGO_TYPE_H_
#define _IMP_ALGO_TYPE_H_


#ifdef __cplusplus
extern "C"
{
#endif

	typedef unsigned char IMP_U8;
	typedef unsigned char IMP_UCHAR;
	typedef unsigned short IMP_U16;
	typedef unsigned int IMP_U32;
	typedef char IMP_S8;
	typedef short IMP_S16;
	typedef int IMP_S32;
#ifndef _M_IX86
	typedef unsigned long long IMP_U64;
	typedef long long IMP_S64;
#else
	typedef __int64 IMP_U64;
	typedef __int64 IMP_S64;
#endif
	typedef char IMP_CHAR;
	typedef char* IMP_PCHAR;
	typedef float IMP_FLOAT;
	typedef double IMP_DOUBLE;
	typedef void IMP_VOID;
	typedef unsigned long IMP_SIZE_T;
	typedef unsigned long IMP_LENGTH_T;

	/* Handle */
	typedef void *IMP_HANDLE;

	typedef enum
	{
		IMP_FALSE = 0,
		IMP_TRUE = 1,
	} IMP_BOOL;

#ifndef NULL
#define NULL 0L
#endif
#define IMP_NULL 0L
#define IMP_NULL_PTR 0L
#define IMP_SUCCESS 0
#define IMP_FAILURE (-1)

	/** IMP_EXPORTS */
#if defined(IMP_API_EXPORTS)
#define IMP_EXPORTS __declspec(dllexport)
#elif defined(IMP_API_IMPORTS)
#define IMP_EXPORTS __declspec(dllimport)
#else
#define IMP_EXPORTS extern
#endif

	/** IMP_INLINE */
#ifndef IMP_INLINE
#if defined __cplusplus
#define IMP_INLINE inline
#elif (defined WIN32 || defined WIN64) && !defined __GNUC__
#define IMP_INLINE __inline
#else
#define IMP_INLINE static
#endif
#endif

	/* Point definition (16b) */
	typedef struct impPOINT_S
	{
		IMP_S16 s16X;
		IMP_S16 s16Y;
	} IMP_POINT_S, IMP_POINT16S_S;

	/* Point definition (32b) */
	typedef struct impPOINT32S_S
	{
		IMP_S32 s32X;
		IMP_S32 s32Y;
	} IMP_POINT32S_S;

	/* Point definition (single-precision) */
	typedef struct impPOINT32F_S
	{
		IMP_FLOAT f32X;
		IMP_FLOAT f32Y;
	} IMP_POINT32F_S;

	/* 3D Point definition (32b) */
	typedef struct impPOINT3D_S
	{
		IMP_S32 s32X;
		IMP_S32 s32Y;
		IMP_S32 s32Z;
	} IMP_POINT3D_S;

	/* Line definition (32b) */
	typedef struct impLINE_S
	{
		IMP_POINT_S stPs; /* start point */
		IMP_POINT_S stPe; /* end point */
	} IMP_LINE_S, LINE_S;

	/* Rectangle definition (32b) */
	typedef struct impIMP_RECT_S
	{
		IMP_S16 s16X1;  /* Top left x */
		IMP_S16 s16Y1;  /* Top left y */
		IMP_S16 s16X2;  /* Bottom right x */
		IMP_S16 s16Y2;  /* Bottom right y */
	} IMP_RECT_S;

	/* Rectangle definition (float)*/
	typedef struct impIMP_RECT32F_S
	{
		IMP_FLOAT f32X1;  /* Top left x */
		IMP_FLOAT f32Y1;  /* Top left y */
		IMP_FLOAT f32X2;  /* Bottom right x */
		IMP_FLOAT f32Y2;  /* Bottom right y */
	} IMP_RECT32F_S;

	/* Polygon definition */
#define MAX_VERTEX_PER_POLYGON 8

	typedef struct impIMP_POLYGON
	{
		IMP_U32  u32VertexNum;    /* Number of vertexes */
		IMP_POINT_S s32Vertex[MAX_VERTEX_PER_POLYGON];  /* Vertexes with clockwise order */
	} IMP_POLYGON_S;

	/* Function Status */
	typedef enum impSTATUS_E
	{
		IMP_STATUS_CHECK_LICENSE_TIMEOUT = -3, /* Timeout error in license checking */
		IMP_STATUS_CHECK_LICENSE_FAILED = -2,  /* Failure in license checking */
		IMP_STATUS_READ_MAC_FAILED = -1,      /* Failed in reading MAC address */
		IMP_STATUS_OK = 1,                    /* Success */
		IMP_STATUS_SKIP,                      /* Function is skipped */
		IMP_STATUS_FALSE,                     /* False in function */
		IMP_STATUS_INVALID_PARA,							/* Function parameters are invalid */
		IMP_STATUS_UNSUPPORT
	} IMP_STATUS_E, STATUS_E;

	/* Pixel definition in RGB */
	typedef struct impPIXEL_S
	{
		IMP_U8 u8B;    /* Blue Channel  */
		IMP_U8 u8G;    /* Green Channel */
		IMP_U8 u8R;    /* Red Channel   */
	} IMP_PIXEL_S;

	/* Pixel definition in HSV */
	typedef struct impHSV_PIXEL_S
	{
		IMP_FLOAT f32H;    /* Hue Channel  */
		IMP_FLOAT f32S;    /* Saturation Channel */
		IMP_FLOAT f32V;    /* Value Channel */
	} IMP_HSV_PIXEL_S;

	/* Pixel definition in HSL */
	typedef struct impHSL_PIXEL_S
	{
		IMP_FLOAT f32H;    /* Hue Channel  */
		IMP_FLOAT f32S;    /* Saturation Channel */
		IMP_FLOAT f32L;    /* Luminance Channel  */
	} IMP_HSL_PIXEL_S;

	/* Image in RGB format */
	typedef struct impRGB_IMAGE_S
	{
		IMP_S32 s32W;      /* Image width  */
		IMP_S32 s32H;      /* Image height */
		IMP_U8 *pu8Data;   /* Image data buffer */
		IMP_U32 u32Time;   /* Time tag */
	} IMP_RGB_IMAGE_S;

	/* Image in HSV format */
	typedef struct impHSV_IMAGE_S
	{
		IMP_S32   s32W;           /* Image width  */
		IMP_S32   s32H;           /* Image height */
		IMP_FLOAT *pf32Data;      /* Image data buffer */
		IMP_U32   u32Time;        /* Time tag */
	} IMP_HSV_IMAGE_S;

	/* Image in gray scale format */
	typedef struct impGRAY_IMAGE_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_U8 *pu8Data;        /* Image data buffer */
	} IMP_GRAY_IMAGE_S;

	/* Image in 16b gray scale format */
	typedef struct impGRAY_IMAGE16_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_S16 *ps16Data;      /* Image data buffer */
	} IMP_GRAY_IMAGE16_S;

	/* Image in 32b gray scale format */
	typedef struct impGRAY_IMAGE32_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_S32 *ps32Data;      /* Image data buffer */
	} IMP_GRAY_IMAGE32_S;

	/* Image in YUV422 format */
	typedef struct impYUV_IMAGE422_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_U8 *pu8Y;           /* Y data buffer */
		IMP_U8 *pu8U;           /* U data buffer */
		IMP_U8 *pu8V;           /* V data buffer */
		IMP_U32 u32Time;        /* Time tag */
		IMP_U32 u32Flag;        /* Reserved flag */
	} IMP_YUV_IMAGE422_S;

	/* Image in YUV420 format */
	typedef struct impYUV_IMAGE420_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_U8 *pu8Y;           /* Y data buffer */
		IMP_U8 *pu8U;           /* U data buffer */
		IMP_U8 *pu8V;           /* V data buffer */
		IMP_U32 u32Time;        /* Time tag */
		IMP_U32 u32Flag;        /* Reserved flag */
	} IMP_YUV_IMAGE420_S;

	/* YUV format type */
	typedef enum impYUV_FORMAT_E
	{
		YUV_FORMAT_IMP_422 = 0,
		YUV_FORMAT_IMP_420,
		YUV_FORMAT_HIS_SP422,
		YUV_FORMAT_HIS_SP420
	} IMP_YUV_FORMAT_E;

	/* Universal YUV image definition */
	typedef struct impYUV_IMAGE_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_U8 *pu8Y;           /* Y data buffer */
		IMP_U8 *pu8U;           /* U data buffer */
		IMP_U8 *pu8V;           /* V data buffer */
		IMP_U32 u32Time;        /* Time tag */
		IMP_U32 u32Flag;        /* Reserved flag */
		IMP_YUV_FORMAT_E enFormat;  /* YUV format */
	} IMP_YUV_IMAGE_S;

	/* 3-ch image format type */
	typedef enum impIMAGE_FORMAT_E
	{
		IMAGE_FORMAT_IMP_YUV422 = 0,
		IMAGE_FORMAT_IMP_YUV420,
		IMAGE_FORMAT_IMP_YUV422SP,
		IMAGE_FORMAT_IMP_YUV420SP,
		IMAGE_FORMAT_IMP_UYVY,
		IMAGE_FORMAT_IMP_YV12,
		IMAGE_FORMAT_IMP_YU12,
		IMAGE_FORMAT_IMP_NV12,
		IMAGE_FORMAT_IMP_NV21,
		IMAGE_FORMAT_IMP_RGB_packed,
		IMAGE_FORMAT_IMP_RGB_planar,
		IMAGE_FORMAT_IMP_BGR_packed,
		IMAGE_FORMAT_IMP_BGR_planar
	} IMP_IMAGE_FORMAT_E;

	/* 3-ch image definition */
	typedef struct impIMAGE3_S
	{
		IMP_S32 s32W;           /* Image width  */
		IMP_S32 s32H;           /* Image height */
		IMP_U8 *pu8D1;          /* Channel 1 data pointer */
		IMP_U8 *pu8D2;          /* Channel 2 data pointer */
		IMP_U8 *pu8D3;          /* Channel 3 data pointer */
		IMP_U32 u32Time;        /* Time tag */
		IMP_U32 u32Flag;        /* Reserved flag */
		IMP_IMAGE_FORMAT_E enFormat;        /* image format pointer */
	} IMP_IMAGE3_S, IMAGE3_S;

#ifdef __cplusplus
}
#endif

#endif /*_IMP_ALGO_TYPE_H_*/


