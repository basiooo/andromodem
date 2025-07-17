import type { MessagesResponse } from "@/types/message"

import {apiClient} from "@/api/apiClient"

export const MessageApi = {
    getMessages: async (serial: string) => {
        const {data} = await apiClient.get<MessagesResponse>(`devices/${serial}/messages`)
        return data
    }
}