import type { AvatarProps } from '@nuxt/ui';

type Paginated<T> = {
  data: T[];
  total: number;
};

export interface User {
  id: number;
  created_at: Date;
  updated_at: Date;
  name: string;
  username: string;
}

export interface Role {
  id: string;
  created_at: Date;
  updated_at: Date;
  name: string;
  description: string;
}

export interface Account {
  id: string;
  created_at: Date;
  updated_at: Date;
  name: string;
  username: string;
  role?: Role;
}

export interface Device {
  id: string;
  created_at: Date;
  updated_at: Date;
  name: string;
  type: string;
  webhook: string;
  jid: string;
  qr_code: string;
  status: string;
  auto_connect: boolean;
  events: string;
  user: User;
}

export interface DeviceToken {
  id: string;
  created_at: Date;
  updated_at: Date;
  token: string;
  status: string;
  user?: User;
}

export interface APIKey {
  id: string;
  created_at: Date;
  updated_at: Date;
  name: string;
  key: string;
  status: string;
  user?: User;
}
