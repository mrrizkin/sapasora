import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
export const create = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: create.url(options),
  method: 'get',
});

create.definition = {
  methods: ['get', 'head'],
  url: '/devices/create',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
create.url = (options?: RouteQueryOptions) => {
  return create.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
create.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
create.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: create.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
export const createForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
createForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: create.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:50
 * @route /devices/create
 */
createForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: create.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
export const edit = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: edit.url(id, options),
  method: 'get',
});

edit.definition = {
  methods: ['get', 'head'],
  url: '/devices/:id/edit',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
edit.url = (id: string, options?: RouteQueryOptions) => {
  return edit.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
edit.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
edit.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: edit.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
export const editForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
editForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: edit.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:89
 * @route /devices/:id/edit
 */
editForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: edit.url(id, options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
export const index = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

index.definition = {
  methods: ['get', 'head'],
  url: '/devices/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
index.url = (options?: RouteQueryOptions) => {
  return index.definition.url + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
index.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
index.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: index.url(options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
export const indexForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
indexForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: index.url(options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:39
 * @route /devices/
 */
indexForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: index.url(options),
  method: 'head',
});
/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
export const show = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: show.url(id, options),
  method: 'get',
});

show.definition = {
  methods: ['get', 'head'],
  url: '/devices/:id/show',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
show.url = (id: string, options?: RouteQueryOptions) => {
  return show.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
show.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
show.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: show.url(id, options),
  method: 'head',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
export const showForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
showForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: show.url(id, options),
  method: 'get',
});

/**
 * @see device/internal/app/http/controllers/device/device.go:61
 * @route /devices/:id/show
 */
showForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: show.url(id, options),
  method: 'head',
});

export const device = {
  create: Object.assign(create, create),
  edit: Object.assign(edit, edit),
  index: Object.assign(index, index),
  show: Object.assign(show, show),
};

export default device;
