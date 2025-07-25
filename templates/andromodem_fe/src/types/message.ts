import type { BaseResponse } from "@/types/response"

export interface Message  {
    row: number;
    address: string;
    body: string;
    date: string;
}
export type Messages = Message[]

export interface MessagesResponse extends BaseResponse {
    data: {
        messages: Messages
    }
}
