import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
export const CheckUser = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: CheckUser.url(options),
  method: 'post',
});

CheckUser.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/check-user',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
CheckUser.url = (options?: RouteQueryOptions) => {
  return CheckUser.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
CheckUser.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: CheckUser.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
export const CheckUserForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: CheckUser.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:46
 * @route /api/v1/gateway/check-user
 */
CheckUserForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: CheckUser.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
export const Connect = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Connect.url(options),
  method: 'post',
});

Connect.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/connect',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
Connect.url = (options?: RouteQueryOptions) => {
  return Connect.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
Connect.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Connect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
export const ConnectForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Connect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:84
 * @route /api/v1/gateway/connect
 */
ConnectForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Connect.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
export const Disconnect = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Disconnect.url(options),
  method: 'post',
});

Disconnect.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/disconnect',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
Disconnect.url = (options?: RouteQueryOptions) => {
  return Disconnect.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
Disconnect.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Disconnect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
export const DisconnectForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Disconnect.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:124
 * @route /api/v1/gateway/disconnect
 */
DisconnectForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Disconnect.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
export const GetAvatar = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetAvatar.url(options),
  method: 'post',
});

GetAvatar.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-avatar',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
GetAvatar.url = (options?: RouteQueryOptions) => {
  return GetAvatar.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
GetAvatar.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetAvatar.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
export const GetAvatarForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetAvatar.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:160
 * @route /api/v1/gateway/get-avatar
 */
GetAvatarForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetAvatar.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
export const GetContacts = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetContacts.url(options),
  method: 'post',
});

GetContacts.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-contacts',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
GetContacts.url = (options?: RouteQueryOptions) => {
  return GetContacts.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
GetContacts.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetContacts.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
export const GetContactsForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetContacts.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:197
 * @route /api/v1/gateway/get-contacts
 */
GetContactsForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetContacts.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
export const GetQR = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetQR.url(options),
  method: 'get',
});

GetQR.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/gateway/get-qr',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
GetQR.url = (options?: RouteQueryOptions) => {
  return GetQR.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
GetQR.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetQR.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
GetQR.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: GetQR.url(options),
  method: 'head',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
export const GetQRForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetQR.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
GetQRForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetQR.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:229
 * @route /api/v1/gateway/get-qr
 */
GetQRForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: GetQR.url(options),
  method: 'head',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
export const GetStatus = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetStatus.url(options),
  method: 'get',
});

GetStatus.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/gateway/get-status',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
GetStatus.url = (options?: RouteQueryOptions) => {
  return GetStatus.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
GetStatus.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetStatus.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
GetStatus.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: GetStatus.url(options),
  method: 'head',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
export const GetStatusForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetStatus.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
GetStatusForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetStatus.url(options),
  method: 'get',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:273
 * @route /api/v1/gateway/get-status
 */
GetStatusForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: GetStatus.url(options),
  method: 'head',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
export const GetUser = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetUser.url(options),
  method: 'post',
});

GetUser.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/get-user',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
GetUser.url = (options?: RouteQueryOptions) => {
  return GetUser.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
GetUser.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: GetUser.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
export const GetUserForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetUser.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:306
 * @route /api/v1/gateway/get-user
 */
GetUserForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: GetUser.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
export const Logout = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Logout.url(options),
  method: 'post',
});

Logout.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/logout',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
Logout.url = (options?: RouteQueryOptions) => {
  return Logout.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
Logout.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Logout.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
export const LogoutForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Logout.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:343
 * @route /api/v1/gateway/logout
 */
LogoutForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Logout.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
export const SendAudio = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendAudio.url(options),
  method: 'post',
});

SendAudio.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-audio',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
SendAudio.url = (options?: RouteQueryOptions) => {
  return SendAudio.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
SendAudio.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendAudio.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
export const SendAudioForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendAudio.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:379
 * @route /api/v1/gateway/send-audio
 */
SendAudioForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendAudio.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
export const SendButton = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendButton.url(options),
  method: 'post',
});

SendButton.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-button',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
SendButton.url = (options?: RouteQueryOptions) => {
  return SendButton.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
SendButton.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendButton.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
export const SendButtonForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendButton.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:417
 * @route /api/v1/gateway/send-button
 */
SendButtonForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendButton.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
export const SendChatPresence = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendChatPresence.url(options),
  method: 'post',
});

SendChatPresence.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-chat-presence',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
SendChatPresence.url = (options?: RouteQueryOptions) => {
  return SendChatPresence.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
SendChatPresence.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendChatPresence.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
export const SendChatPresenceForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendChatPresence.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:455
 * @route /api/v1/gateway/send-chat-presence
 */
SendChatPresenceForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendChatPresence.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
export const SendContact = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendContact.url(options),
  method: 'post',
});

SendContact.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-contact',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
SendContact.url = (options?: RouteQueryOptions) => {
  return SendContact.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
SendContact.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendContact.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
export const SendContactForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendContact.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:496
 * @route /api/v1/gateway/send-contact
 */
SendContactForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendContact.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
export const SendDocument = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendDocument.url(options),
  method: 'post',
});

SendDocument.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-document',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
SendDocument.url = (options?: RouteQueryOptions) => {
  return SendDocument.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
SendDocument.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendDocument.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
export const SendDocumentForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendDocument.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:534
 * @route /api/v1/gateway/send-document
 */
SendDocumentForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendDocument.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
export const SendImage = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendImage.url(options),
  method: 'post',
});

SendImage.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-image',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
SendImage.url = (options?: RouteQueryOptions) => {
  return SendImage.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
SendImage.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendImage.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
export const SendImageForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendImage.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:572
 * @route /api/v1/gateway/send-image
 */
SendImageForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendImage.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
export const SendList = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendList.url(options),
  method: 'post',
});

SendList.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-list',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
SendList.url = (options?: RouteQueryOptions) => {
  return SendList.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
SendList.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendList.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
export const SendListForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendList.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:610
 * @route /api/v1/gateway/send-list
 */
SendListForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendList.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
export const SendLocation = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendLocation.url(options),
  method: 'post',
});

SendLocation.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-location',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
SendLocation.url = (options?: RouteQueryOptions) => {
  return SendLocation.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
SendLocation.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendLocation.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
export const SendLocationForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendLocation.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:648
 * @route /api/v1/gateway/send-location
 */
SendLocationForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendLocation.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
export const SendSticker = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendSticker.url(options),
  method: 'post',
});

SendSticker.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-sticker',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
SendSticker.url = (options?: RouteQueryOptions) => {
  return SendSticker.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
SendSticker.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendSticker.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
export const SendStickerForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendSticker.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:686
 * @route /api/v1/gateway/send-sticker
 */
SendStickerForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendSticker.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
export const SendText = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendText.url(options),
  method: 'post',
});

SendText.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-text',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
SendText.url = (options?: RouteQueryOptions) => {
  return SendText.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
SendText.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendText.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
export const SendTextForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendText.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:724
 * @route /api/v1/gateway/send-text
 */
SendTextForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendText.url(options),
  method: 'post',
});
/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
export const SendVideo = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendVideo.url(options),
  method: 'post',
});

SendVideo.definition = {
  methods: ['post'],
  url: '/api/v1/gateway/send-video',
} satisfies RouteDefinition<['post']>;

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
SendVideo.url = (options?: RouteQueryOptions) => {
  return SendVideo.definition.url + queryParams(options);
};

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
SendVideo.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: SendVideo.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
export const SendVideoForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendVideo.url(options),
  method: 'post',
});

/**
 * @see gateway/internal/app/http/controllers/gateway/gateway.go:762
 * @route /api/v1/gateway/send-video
 */
SendVideoForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: SendVideo.url(options),
  method: 'post',
});

export const GatewayController = {
  CheckUser: Object.assign(CheckUser, CheckUser),
  Connect: Object.assign(Connect, Connect),
  Disconnect: Object.assign(Disconnect, Disconnect),
  GetAvatar: Object.assign(GetAvatar, GetAvatar),
  GetContacts: Object.assign(GetContacts, GetContacts),
  GetQR: Object.assign(GetQR, GetQR),
  GetStatus: Object.assign(GetStatus, GetStatus),
  GetUser: Object.assign(GetUser, GetUser),
  Logout: Object.assign(Logout, Logout),
  SendAudio: Object.assign(SendAudio, SendAudio),
  SendButton: Object.assign(SendButton, SendButton),
  SendChatPresence: Object.assign(SendChatPresence, SendChatPresence),
  SendContact: Object.assign(SendContact, SendContact),
  SendDocument: Object.assign(SendDocument, SendDocument),
  SendImage: Object.assign(SendImage, SendImage),
  SendList: Object.assign(SendList, SendList),
  SendLocation: Object.assign(SendLocation, SendLocation),
  SendSticker: Object.assign(SendSticker, SendSticker),
  SendText: Object.assign(SendText, SendText),
  SendVideo: Object.assign(SendVideo, SendVideo),
};

export default GatewayController;
