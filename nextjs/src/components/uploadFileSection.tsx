import React, { useRef, useState } from "react";
import UploadButton from "./uploadButton";
import { getBase64FromImageFile } from "@/utils/processImage";
import { ClipLoader } from "react-spinners";

interface FileSectionProps {
  // eslint-disable-next-line no-unused-vars
  setDisplay: (image: string, caption: string) => void;
  error: string | null;
  // eslint-disable-next-line no-unused-vars
  setError: (error: string) => void;
}

const ALLOWED_TYPES = ["image/png", "image/jpeg", "image/webp", "image/heic"];

const UploadFileSection: React.FC<FileSectionProps> = ({
  setDisplay,
  error,
  setError,
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);

  const handleComputerClick = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const isValidFile = (file: File) => {
    return ALLOWED_TYPES.includes(file.type);
  };

  const handleFileChange = async (file: File) => {
    // Check file type
    if (!isValidFile(file)) {
      setError("Invalid file type - upload a PNG, JPEG, WEBP, or HEIC image.");
      return;
    }

    // Check file size
    if (file.size > 5 * 1024 * 1024) {
      setError("Invalid file size - must be less than 5 MB.");
      return;
    }

    setLoading(true);

    try {
      const base64String = await getBase64FromImageFile(file); // data URL
      const base64Data = base64String.split(",")[1]; // just base64 encoding

      console.log("Calling caption api");
      const response = await fetch("/api/caption", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ base64encoding: base64Data }),
      });

      let jobID = "";

      if (response.ok) {
        const responseBody = await response.json();
        jobID = responseBody["apiData"]["body"];
      }

      const statusResponse = await fetch(`/api/job/status/${jobID}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      let caption = "";

      if (statusResponse.ok) {
        const responseBody = await statusResponse.json();
        const status = responseBody["apiData"]["Status"];
        if (status === "Failed") {
          setError("Something went wrong - try another image.");
        } else {
          caption = responseBody["apiData"]["Caption"];
          setDisplay(base64String, caption);
        }
      } else {
        setError("Something went wrong - try another image.");
      }
    } catch (error) {
      setError("Something went wrong - try another image.");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const file = event.dataTransfer.files?.[0];
    if (file) {
      handleFileChange(file);
    }
  };

  const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
  };

  return (
    <div
      className="flex min-w-[300px] flex-col items-center justify-center rounded-lg bg-white p-2 pb-10 pt-20 text-xl shadow-lg md:px-40 lg:text-3xl"
      onDrop={handleDrop}
      onDragOver={handleDragOver}
    >
      <div className="flex h-48 w-full max-w-md items-center justify-center">
        {loading ? (
          <div className="flex h-full w-full flex-col items-center justify-center gap-6 px-8">
            <ClipLoader color="#3498db" />
            <div className="max-w-60 text-center text-base text-gray-400 md:max-w-80">
              Generating caption - do not exit page.
            </div>
          </div>
        ) : (
          <div className="flex h-full w-full flex-col items-center justify-center">
            <div className="mb-4 flex items-center text-center font-bold text-gray-400">
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
                className="mr-2 text-gray-400"
              >
                <g id="SVGRepo_bgCarrier" strokeWidth="0"></g>
                <g
                  id="SVGRepo_tracerCarrier"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                ></g>
                <g id="SVGRepo_iconCarrier">
                  <path
                    fillRule="evenodd"
                    clipRule="evenodd"
                    d="M23 4C23 2.34315 21.6569 1 20 1H4C2.34315 1 1 2.34315 1 4V20C1 21.6569 2.34315 23 4 23H20C21.6569 23 23 21.6569 23 20V4ZM21 4C21 3.44772 20.5523 3 20 3H4C3.44772 3 3 3.44772 3 4V20C3 20.5523 3.44772 21 4 21H20C20.5523 21 21 20.5523 21 20V4Z"
                    fill="currentColor"
                  ></path>
                  <path
                    d="M4.80665 17.5211L9.1221 9.60947C9.50112 8.91461 10.4989 8.91461 10.8779 9.60947L14.0465 15.4186L15.1318 13.5194C15.5157 12.8476 16.4843 12.8476 16.8682 13.5194L19.1451 17.5039C19.526 18.1705 19.0446 19 18.2768 19H5.68454C4.92548 19 4.44317 18.1875 4.80665 17.5211Z"
                    fill="currentColor"
                  ></path>
                  <path
                    d="M18 8C18 9.10457 17.1046 10 16 10C14.8954 10 14 9.10457 14 8C14 6.89543 14.8954 6 16 6C17.1046 6 18 6.89543 18 8Z"
                    fill="currentColor"
                  ></path>
                </g>
              </svg>
              <div className="hidden text-center md:block">
                Drag your photo here!
              </div>
              <div className="block text-center md:hidden">
                Select image here!
              </div>
            </div>

            <div className="mb-4 flex w-full items-center text-xl font-semibold text-gray-400">
              <div className="mr-4 flex-1 border-t border-gray-400 text-center"></div>
              <div className="hidden text-center md:block">or upload from</div>
              <div className="block text-center md:hidden">upload from</div>
              <div className="ml-4 flex-1 border-t border-gray-400"></div>
            </div>
            <div className="pt-5">
              <UploadButton text="Device" onClick={handleComputerClick} />
            </div>
            <div className="max-w-60 py-5 text-center text-sm text-red-500 md:max-w-80">
              {error ? error : ""}
            </div>
          </div>
        )}
      </div>
      <input
        type="file"
        ref={fileInputRef}
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) {
            handleFileChange(file);
          }
        }}
        className="hidden"
        accept=".png, .jpg, .jpeg, .HEIC, .webp"
      />
    </div>
  );
};

export default UploadFileSection;
