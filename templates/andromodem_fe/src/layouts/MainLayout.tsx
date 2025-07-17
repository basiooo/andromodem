import type { FC } from "react"
import { useEffect } from "react"

import Content from "@/components/Content/Content"
import DeviceSelector from "@/components/DeviceSelector/DeviceSelector"
import Footer from "@/components/Footer/Footer"
import Header from "@/components/Header/Header"
import UpdateModal from "@/components/Modal/UpdateModal"
import useHealthCheck from "@/hooks/useHealthCheck"
import useUpdateChecker from "@/hooks/useUpdateChecker"
import { useUpdateStore } from "@/stores/updateStore"

const MainLayout: FC = () => {
  const { updateInfo: checkerUpdateInfo } = useUpdateChecker()
  const { updateInfo, setUpdateInfo } = useUpdateStore()
  const {error: healthError } = useHealthCheck()
  
  if (healthError) {
    throw healthError
  }

  useEffect(() => {
    if (checkerUpdateInfo) {
      setUpdateInfo(checkerUpdateInfo)
    }
  }, [checkerUpdateInfo, setUpdateInfo])

  return (
    <>
      <Header/>
      <main className="my-5 md:my-15">
        <div className="container mx-auto">
          <DeviceSelector/>
          <Content/>
        </div>
      </main>
      <Footer/>
      
      {/* Update Modal */}
      <UpdateModal updateInfo={updateInfo} modal_id="updateModal" />
    </>
  )
}

export default MainLayout