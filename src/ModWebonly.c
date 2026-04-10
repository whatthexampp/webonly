#ifdef _WIN32
#include <winsock2.h>
#include <ws2tcpip.h>
#pragma comment(lib, "ws2_32.lib")
#define CLOSE_SOCKET closesocket
#else
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#define CLOSE_SOCKET close
#endif

#include "httpd.h"
#include "http_config.h"
#include "http_protocol.h"
#include "ap_config.h"

static int WebonlyHandler(request_rec *Req) {
    if (!Req->handler || strcmp(Req->handler, "webonly-script") != 0) {
        return DECLINED;
    }

    ap_set_content_type(Req, "text/html");

#ifdef _WIN32
    WSADATA Wsa;
    WSAStartup(MAKEWORD(2,2), &Wsa);
#endif

    int Sock = socket(AF_INET, SOCK_STREAM, 0);
    if (Sock < 0) {
        ap_rputs("socket error", Req);
        return OK;
    }

    struct sockaddr_in Addr;
    memset(&Addr, 0, sizeof(Addr));
    Addr.sin_family = AF_INET;
    Addr.sin_port = htons(43211);
    inet_pton(AF_INET, "127.0.0.1", &Addr.sin_addr);

    if (connect(Sock, (struct sockaddr *)&Addr, sizeof(Addr)) < 0) {
        ap_rputs("run daemon:", Req);
        CLOSE_SOCKET(Sock);
        return OK;
    }

    send(Sock, Req->filename, strlen(Req->filename), 0);
    send(Sock, "?", 1, 0);
    if (Req->args) {
        send(Sock, Req->args, strlen(Req->args), 0);
    }
    send(Sock, "\n", 1, 0);

    char Buf[4096];
    int Read;
    while ((Read = recv(Sock, Buf, sizeof(Buf) - 1, 0)) > 0) {
        Buf[Read] = '\0';
        ap_rputs(Buf, Req);
    }

    CLOSE_SOCKET(Sock);

#ifdef _WIN32
    WSACleanup();
#endif

    return OK;
}

static void RegisterWebonlyHooks(apr_pool_t *Pool) {
    ap_hook_handler(WebonlyHandler, NULL, NULL, APR_HOOK_MIDDLE);
}

module AP_MODULE_DECLARE_DATA WebonlyModule = {
    STANDARD20_MODULE_STUFF,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    RegisterWebonlyHooks
};