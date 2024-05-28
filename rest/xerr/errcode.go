package xerr

// 成功返回
const OK uint32 = 1000

// 失败
const SERVER_COMMON_ERROR uint32 = 1001

// 参数有误
const REUQEST_PARAM_ERROR uint32 = 1002

// 登录有误
const LOGIN_EXPIRE_ERROR uint32 = 1003

// 验签错误
const SIGN_ERROR uint32 = 1004

// 没有数据权限错误
const MISSED_DATA_PERMISSIONS_ERROR uint32 = 1005

// http服务请求错误
const HTTP_SERVER_REQUEST_ERROR uint32 = 1006

// 第三方服务错误
const THIRD_SERVER_REQUEST_ERROR uint32 = 1007

// 没有功能权限错误
const MISSED_FUNC_PERMISSIONS_ERROR uint32 = 1008
