import { z } from "zod";

export const Id = z.string().cuid2();
