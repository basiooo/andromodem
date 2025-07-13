import { create } from 'zustand'

import type { UpdateStore } from '@/types/update'

export const useUpdateStore = create<UpdateStore>((set) => {
  return{
    updateInfo: null,
    setUpdateInfo: (info) => set({ updateInfo: info })
  }
})