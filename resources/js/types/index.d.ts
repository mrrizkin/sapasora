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
  events: string;
  user: User;
}
