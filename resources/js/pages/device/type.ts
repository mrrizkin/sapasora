import * as z from 'zod';

export type ColumnFilter = {
  search?: string;
};

export const deviceFormSchema = z.object({
  name: z.string().min(1, 'Device name is required'),
  type: z.string().min(1, 'Device type is required'),
  webhook: z.url('Invalid URL').optional().or(z.literal('')),
});

export type DeviceFormSchema = z.output<typeof deviceFormSchema>;
