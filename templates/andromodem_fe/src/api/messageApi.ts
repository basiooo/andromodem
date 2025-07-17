import {apiClient} from "@/api/apiClient"
import type { MessagesResponse } from "@/types/message"

export const MessageApi = {
    getMessages: async (serial: string) => {
        const {data} = await apiClient.get<MessagesResponse>(`devices/${serial}/messages`)
        return data
    }
}