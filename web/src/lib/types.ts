import { z } from "zod";

export const SetupSchema = z.object({
  user: z
    .object({
      username: z.string(),
      password: z.string(),
    })
    .optional(),
});

export type Setup = z.infer<typeof SetupSchema>;
