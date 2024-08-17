"use client";
import Example from "@/components/example";
import UploadButton from "@/components/uploadButton";
import UploadFileSection from "@/components/uploadFileSection";
import shirt from "@/assets/shirt.jpeg";
import DisplayImageSection from "@/components/displayImageSection";

export default function Home() {
  return (
    <main className="min-h-[100vh] overflow-hidden bg-white">
      {/* Header */}
      <div className="h-16 w-full bg-white"></div>
      {/* Upload image section */}
      <div className="w-full bg-[#668F80] p-10">
        <div className="pb-5">
          <div className="flex items-center justify-center py-2 text-center text-3xl font-semibold">
            CAPTION IMAGE
          </div>
          <div className="mb-4 flex items-center justify-center text-center text-xl font-semibold">
            Create a title for your clothing listing!
          </div>
        </div>
        <div className="flex items-center justify-center font-semibold">
          {/* <UploadFileSection /> */}
          <DisplayImageSection title="test" image={shirt} />
        </div>
      </div>
      {/* Examples section */}
      <div className="flex flex-col items-center justify-center p-7">
        <div className="flex items-center justify-center pb-2 text-2xl text-gray-500 md:w-[500px] md:text-3xl">
          <div className="mr-4 flex-1 border-t border-gray-400"></div>
          Examples
          <div className="ml-4 flex-1 border-t border-gray-400"></div>
        </div>
        <div className="mt-3 flex flex-col items-center justify-center gap-10 px-10 md:flex-row md:items-start md:gap-20">
          <div className="h-[250px] w-[250px] md:h-[200px] md:w-[200px]">
            <Example title="dsadsa" image={shirt} />
          </div>
          <div className="h-[250px] w-[250px] md:h-[200px] md:w-[200px]">
            <Example title="dsadsa" image={shirt} />
          </div>
          <div className="h-[250px] w-[250px] md:h-[200px] md:w-[200px]">
            <Example title="dsadsa" image={shirt} />
          </div>
        </div>
      </div>
    </main>
  );
}
