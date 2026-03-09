import type { UserType } from "@/types/user";

// 泛型，定義共通 Response
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface ResponseType<T = any> {
  message: string;
  data: T;
  timeStamp?: string;
}

type FileResponseData = {
  files: string[];
  directories: string[];
};

export type FileResponseType = ResponseType<FileResponseData>;
export type LoginResponseType = ResponseType<UserType>;
