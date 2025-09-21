import { type FC } from "react"

import type { Device } from "@/types/device"

import MirroringCanvas from './MirroringCanvas'

const Mirroring: FC<{ device: Device }> = ({
  device
}) => {
  return (
    <MirroringCanvas
        device={device}
      />
  )
}
export default Mirroring