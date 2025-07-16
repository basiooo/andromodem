import type { FC } from "react"
import { useState } from "react"
import { FaDownload } from "react-icons/fa"
import { MdFullscreen, MdFullscreenExit } from "react-icons/md"

import ThemeSelector from "@/components/ThemeSelector/ThemeSelector"
import { useUpdateStore } from "@/stores/updateStore"
import { showModal } from "@/utils/common"

const Header: FC = () => {
  const [isFullscreen, setIsFullscreen] = useState(false)
  const updateInfo = useUpdateStore((state) => state.updateInfo)

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen()
        .then(() => setIsFullscreen(true))
        .catch(() => {})
    } else {
      document.exitFullscreen()
        .then(() => setIsFullscreen(false))
        .catch(() => {})
    }
  }

  return (
    <div className="navbar sticky min-h-20 top-0 backdrop-blur-sm bg-base-300 z-40 shadow-lg">
      <div className="container flex mx-auto">
        <div className="flex-1">
          <a className="text-xl md:text-3xl">AndroModem</a>
        </div>
        <div className="flex gap-2 items-center">
          <ThemeSelector />
          <button
            onClick={toggleFullscreen}
            aria-label={isFullscreen ? "Exit fullscreen" : "Enter fullscreen"}
            className="btn btn-square btn-ghost text-2xl"
            title={isFullscreen ? "Exit Fullscreen" : "Enter Fullscreen"}
          >
            {isFullscreen ? <MdFullscreenExit /> : <MdFullscreen />}
          </button>
          {updateInfo?.hasUpdate && (
            <button
              onClick={() => showModal("updateModal")}
              className="btn btn-primary btn-xs gap-2"
              title="Update tersedia"
            >
              <FaDownload className="text-sm" />
              Update
            </button>
          )}
        </div>
      </div>
    </div>
  )
}

export default Header
