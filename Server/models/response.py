from pydantic import BaseModel
from typing import Optional, Any, Generic, TypeVar

T = TypeVar("T")

# 成功状态码
SUCCESS_CODE = 20000
# 错误状态码
ERROR_CODE = 50000
ILLEGAL_TOKEN_CODE = 50008
OTHER_CLIENT_LOGIN_CODE = 50012
TOKEN_EXPIRED_CODE = 50014


class BaseResponse(BaseModel, Generic[T]):
    code: int
    status: str
    msg: str
    data: Optional[T] = None


def success_response(data: Any = None, msg: str = "成功") -> BaseResponse:
    """成功响应"""
    return BaseResponse(code=SUCCESS_CODE, status="ok", msg=msg, data=data)


def error_response(code: int = ERROR_CODE, msg: str = "错误") -> BaseResponse:
    """错误响应"""
    return BaseResponse(code=code, status="fail", msg=msg, data=None)


def invalid_credentials_response() -> BaseResponse:
    """无效凭证错误"""
    return error_response(ERROR_CODE, "用户名或密码错误")


def illegal_token_response() -> BaseResponse:
    """非法token错误"""
    return error_response(ILLEGAL_TOKEN_CODE, "非法token")


def other_client_login_response() -> BaseResponse:
    """其他客户端已登录"""
    return error_response(OTHER_CLIENT_LOGIN_CODE, "其他客户端已登录")


def token_expired_response() -> BaseResponse:
    """Token已过期"""
    return error_response(TOKEN_EXPIRED_CODE, "Token已过期")
