import * as z from 'zod';

export type ColumnFilter = {
  search?: string;
};

export const userFormSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
  password: z.string().min(3, 'Password must be at least 3 characters'),
  role_id: z.string().min(1, 'Role is required'),
});

export type UserFormSchema = z.output<typeof userFormSchema>;
