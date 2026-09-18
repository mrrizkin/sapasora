import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
export const check_user = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: check_user.url(options),
  method: 'post',
});

check_user.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/check-user',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
check_user.url = (options?: RouteQueryOptions) => {
  return check_user.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
check_user.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: check_user.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
export const check_userForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: check_user.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
check_userForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: check_user.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
export const connect = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: connect.url(options),
  method: 'post',
});

connect.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/connect',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
connect.url = (options?: RouteQueryOptions) => {
  return connect.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
connect.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: connect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
export const connectForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: connect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
connectForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: connect.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
export const disconnect = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: disconnect.url(options),
  method: 'post',
});

disconnect.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/disconnect',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
disconnect.url = (options?: RouteQueryOptions) => {
  return disconnect.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
disconnect.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: disconnect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
export const disconnectForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: disconnect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
disconnectForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: disconnect.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
export const get_avatar = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_avatar.url(options),
  method: 'post',
});

get_avatar.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-avatar',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
get_avatar.url = (options?: RouteQueryOptions) => {
  return get_avatar.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
get_avatar.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_avatar.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
export const get_avatarForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_avatar.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
get_avatarForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_avatar.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
export const get_contacts = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_contacts.url(options),
  method: 'post',
});

get_contacts.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-contacts',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
get_contacts.url = (options?: RouteQueryOptions) => {
  return get_contacts.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
get_contacts.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_contacts.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
export const get_contactsForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_contacts.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
get_contactsForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_contacts.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
export const get_qr = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_qr.url(options),
  method: 'get',
});

get_qr.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/gateway/get-qr',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
get_qr.url = (options?: RouteQueryOptions) => {
  return get_qr.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
get_qr.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_qr.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
get_qr.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get_qr.url(options),
  method: 'head',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
export const get_qrForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_qr.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
get_qrForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_qr.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
get_qrForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get_qr.url(options),
  method: 'head',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
export const get_status = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_status.url(options),
  method: 'get',
});

get_status.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/gateway/get-status',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
get_status.url = (options?: RouteQueryOptions) => {
  return get_status.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
get_status.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get_status.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
get_status.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get_status.url(options),
  method: 'head',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
export const get_statusForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_status.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
get_statusForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get_status.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
get_statusForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get_status.url(options),
  method: 'head',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
export const get_user = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_user.url(options),
  method: 'post',
});

get_user.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-user',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
get_user.url = (options?: RouteQueryOptions) => {
  return get_user.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
get_user.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: get_user.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
export const get_userForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_user.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
get_userForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: get_user.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
export const logout = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: logout.url(options),
  method: 'post',
});

logout.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/logout',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
logout.url = (options?: RouteQueryOptions) => {
  return logout.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
logout.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: logout.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
export const logoutForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: logout.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
logoutForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: logout.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
export const send_audio = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_audio.url(options),
  method: 'post',
});

send_audio.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-audio',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
send_audio.url = (options?: RouteQueryOptions) => {
  return send_audio.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
send_audio.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_audio.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
export const send_audioForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_audio.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
send_audioForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_audio.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
export const send_button = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_button.url(options),
  method: 'post',
});

send_button.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-button',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
send_button.url = (options?: RouteQueryOptions) => {
  return send_button.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
send_button.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_button.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
export const send_buttonForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_button.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
send_buttonForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_button.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
export const send_chat_presence = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_chat_presence.url(options),
  method: 'post',
});

send_chat_presence.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-chat-presence',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
send_chat_presence.url = (options?: RouteQueryOptions) => {
  return send_chat_presence.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
send_chat_presence.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_chat_presence.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
export const send_chat_presenceForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_chat_presence.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
send_chat_presenceForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_chat_presence.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
export const send_contact = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_contact.url(options),
  method: 'post',
});

send_contact.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-contact',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
send_contact.url = (options?: RouteQueryOptions) => {
  return send_contact.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
send_contact.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_contact.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
export const send_contactForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_contact.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
send_contactForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_contact.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
export const send_document = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_document.url(options),
  method: 'post',
});

send_document.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-document',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
send_document.url = (options?: RouteQueryOptions) => {
  return send_document.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
send_document.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_document.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
export const send_documentForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_document.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
send_documentForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_document.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
export const send_image = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_image.url(options),
  method: 'post',
});

send_image.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-image',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
send_image.url = (options?: RouteQueryOptions) => {
  return send_image.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
send_image.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_image.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
export const send_imageForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_image.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
send_imageForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_image.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
export const send_list = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_list.url(options),
  method: 'post',
});

send_list.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-list',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
send_list.url = (options?: RouteQueryOptions) => {
  return send_list.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
send_list.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_list.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
export const send_listForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_list.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
send_listForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_list.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
export const send_location = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_location.url(options),
  method: 'post',
});

send_location.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-location',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
send_location.url = (options?: RouteQueryOptions) => {
  return send_location.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
send_location.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_location.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
export const send_locationForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_location.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
send_locationForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_location.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
export const send_sticker = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_sticker.url(options),
  method: 'post',
});

send_sticker.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-sticker',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
send_sticker.url = (options?: RouteQueryOptions) => {
  return send_sticker.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
send_sticker.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_sticker.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
export const send_stickerForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_sticker.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
send_stickerForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_sticker.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
export const send_text = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_text.url(options),
  method: 'post',
});

send_text.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-text',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
send_text.url = (options?: RouteQueryOptions) => {
  return send_text.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
send_text.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_text.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
export const send_textForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_text.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
send_textForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_text.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
export const send_video = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_video.url(options),
  method: 'post',
});

send_video.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-video',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
send_video.url = (options?: RouteQueryOptions) => {
  return send_video.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
send_video.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: send_video.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
export const send_videoForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_video.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
send_videoForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: send_video.url(options),
  method: 'post',
});

export const api_v1_gateway = {
  check_user: Object.assign(check_user, check_user),
  connect: Object.assign(connect, connect),
  disconnect: Object.assign(disconnect, disconnect),
  get_avatar: Object.assign(get_avatar, get_avatar),
  get_contacts: Object.assign(get_contacts, get_contacts),
  get_qr: Object.assign(get_qr, get_qr),
  get_status: Object.assign(get_status, get_status),
  get_user: Object.assign(get_user, get_user),
  logout: Object.assign(logout, logout),
  send_audio: Object.assign(send_audio, send_audio),
  send_button: Object.assign(send_button, send_button),
  send_chat_presence: Object.assign(send_chat_presence, send_chat_presence),
  send_contact: Object.assign(send_contact, send_contact),
  send_document: Object.assign(send_document, send_document),
  send_image: Object.assign(send_image, send_image),
  send_list: Object.assign(send_list, send_list),
  send_location: Object.assign(send_location, send_location),
  send_sticker: Object.assign(send_sticker, send_sticker),
  send_text: Object.assign(send_text, send_text),
  send_video: Object.assign(send_video, send_video),
};

export default api_v1_gateway;
