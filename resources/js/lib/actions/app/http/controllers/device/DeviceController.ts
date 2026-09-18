import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
export const Create = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Create.url(options),
  method: 'get',
});

Create.definition = {
  methods: ['get', 'head'],
  url: '/devices/create',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
Create.url = (options?: RouteQueryOptions) => {
  return Create.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
Create.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
Create.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Create.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
export const CreateForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
CreateForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
CreateForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Create.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
export const Destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

Destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
Destroy.url = (id: string, options?: RouteQueryOptions) => {
  return Destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
Destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
export const DestroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:418
 * @route /api/v1/device/:id
 */
DestroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
export const Edit = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Edit.url(id, options),
  method: 'get',
});

Edit.definition = {
  methods: ['get', 'head'],
  url: '/devices/:id/edit',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
Edit.url = (id: string, options?: RouteQueryOptions) => {
  return Edit.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
Edit.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
Edit.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Edit.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
export const EditForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
EditForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
EditForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Edit.url(id, options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
export const Get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

Get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
Get.url = (id: string, options?: RouteQueryOptions) => {
  return Get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
Get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
Get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Get.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
export const GetForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
GetForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:150
 * @route /api/v1/device/:id
 */
GetForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Get.url(id, options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
export const GetDeviceByToken = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetDeviceByToken.url(options),
  method: 'get',
});

GetDeviceByToken.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/by-token',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
GetDeviceByToken.url = (options?: RouteQueryOptions) => {
  return GetDeviceByToken.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
GetDeviceByToken.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: GetDeviceByToken.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
GetDeviceByToken.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: GetDeviceByToken.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
export const GetDeviceByTokenForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetDeviceByToken.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
GetDeviceByTokenForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: GetDeviceByToken.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:200
 * @route /api/v1/device/by-token
 */
GetDeviceByTokenForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: GetDeviceByToken.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
export const Index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

Index.definition = {
  methods: ['get', 'head'],
  url: '/devices/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
Index.url = (options?: RouteQueryOptions) => {
  return Index.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
Index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
Index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Index.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
export const IndexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
IndexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
IndexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Index.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
export const List = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

List.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/device/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
List.url = (options?: RouteQueryOptions) => {
  return List.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
List.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
List.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: List.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
export const ListForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
ListForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:118
 * @route /api/v1/device/
 */
ListForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: List.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
export const Show = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Show.url(id, options),
  method: 'get',
});

Show.definition = {
  methods: ['get', 'head'],
  url: '/devices/:id/show',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
Show.url = (id: string, options?: RouteQueryOptions) => {
  return Show.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
Show.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
Show.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Show.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
export const ShowForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
ShowForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
ShowForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Show.url(id, options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
export const Store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

Store.definition = {
  methods: ['post'],
  url: '/api/v1/device/',
} satisfies RouteDefinition<['post']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
Store.url = (options?: RouteQueryOptions) => {
  return Store.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
Store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
export const StoreForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:210
 * @route /api/v1/device/
 */
StoreForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
export const Update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

Update.definition = {
  methods: ['put'],
  url: '/api/v1/device/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
Update.url = (id: string, options?: RouteQueryOptions) => {
  return Update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
Update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
export const UpdateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:292
 * @route /api/v1/device/:id
 */
UpdateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
export const UpdateStatus = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdateStatus.url(id, options),
  method: 'put',
});

UpdateStatus.definition = {
  methods: ['put'],
  url: '/api/v1/device/:id/status',
} satisfies RouteDefinition<['put']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
UpdateStatus.url = (id: string, options?: RouteQueryOptions) => {
  return UpdateStatus.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
UpdateStatus.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdateStatus.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
export const UpdateStatusForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdateStatus.url(id, options),
  method: 'put',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:370
 * @route /api/v1/device/:id/status
 */
UpdateStatusForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdateStatus.url(id, options),
  method: 'put',
});

export const DeviceController = {
  Create: Object.assign(Create, Create),
  Destroy: Object.assign(Destroy, Destroy),
  Edit: Object.assign(Edit, Edit),
  Get: Object.assign(Get, Get),
  GetDeviceByToken: Object.assign(GetDeviceByToken, GetDeviceByToken),
  Index: Object.assign(Index, Index),
  List: Object.assign(List, List),
  Show: Object.assign(Show, Show),
  Store: Object.assign(Store, Store),
  Update: Object.assign(Update, Update),
  UpdateStatus: Object.assign(UpdateStatus, UpdateStatus),
};

export default DeviceController;
