#ifndef MPV_CLIENT_H_
#define MPV_CLIENT_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct mpv_handle mpv_handle;

typedef enum mpv_event_id {
    MPV_EVENT_NONE = 0,
    MPV_EVENT_SHUTDOWN = 1,
    MPV_EVENT_LOG_MESSAGE = 2,
    MPV_EVENT_GET_PROPERTY_REPLY = 3,
    MPV_EVENT_SET_PROPERTY_REPLY = 4,
    MPV_EVENT_COMMAND_REPLY = 5,
    MPV_EVENT_START_FILE = 6,
    MPV_EVENT_END_FILE = 7,
    MPV_EVENT_FILE_LOADED = 8,
    MPV_EVENT_PLAYBACK_RESTART = 21,
} mpv_event_id;

typedef struct mpv_event {
    mpv_event_id event_id;
    int error;
    uint64_t reply_userdata;
    void *data;
} mpv_event;

// MPV API functions
mpv_handle *mpv_create(void);
int mpv_initialize(mpv_handle *ctx);
void mpv_terminate_destroy(mpv_handle *ctx);
int mpv_set_option_string(mpv_handle *ctx, const char *name, const char *data);
int mpv_command(mpv_handle *ctx, const char **args);
mpv_event *mpv_wait_event(mpv_handle *ctx, double timeout);
int mpv_request_log_messages(mpv_handle *ctx, const char *min_level);
int mpv_set_property_string(mpv_handle *ctx, const char *name, const char *data);

#ifdef __cplusplus
}
#endif

#endif
