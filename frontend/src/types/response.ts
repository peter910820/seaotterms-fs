// 泛型，定義共通 Response
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface ResponseType<T = any> {
  message: string;
  data: T;
  timeStamp?: string;
}

export interface FileResponseData {
  files: string[];
  directories: string[];
}

export interface LoginResponseData {
  username: string;
  email: string;
  avatar: string;
  isAdmin: boolean;
  createdAt: string;
}
