import React, { useState } from "react";
import UploadButton from "./uploadButton";
import { StaticImageData } from "next/image";

interface DisplaySectionProps {
  title: string;
  image: string | StaticImageData;
}

const DisplayImageSection: React.FC<DisplaySectionProps> = ({
  title,
  image,
}) => {
  const [text, setText] = useState(title);

  const handleTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    setText(event.target.value);
  };

  return (
    <div className="flex min-w-[300px] flex-col items-center justify-center rounded-lg bg-white p-2 text-xl shadow-lg lg:text-3xl">
      <div className="w-full">
        <button className="flex items-center rounded p-2 text-gray-500 hover:text-gray-300">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2}
            stroke="currentColor"
            className="h-5 w-5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M10 19l-7-7m0 0l7-7m-7 7h18"
            />
          </svg>{" "}
          <div className="px-2 text-lg">Back</div>
        </button>
      </div>
      <div className="flex flex-col items-center justify-center gap-5 py-10 md:flex-row md:px-20">
        <div className="pt-5 shadow-lg md:mr-10">
          <img
            src={typeof image === "string" ? image : image.src}
            alt={title}
            className="max-h-64 max-w-64"
          />
        </div>
        <div className="flex flex-col items-end gap-2">
          <textarea
            value={text}
            onChange={handleTextChange}
            className="resiz-none max-h-6 rounded border border-gray-300 p-2 text-xl text-gray-500"
            rows={1} // Set initial number of rows
            style={{ minHeight: "4rem" }} // Ensure a minimum height
          />
          <button className="h-8 w-20 rounded-lg bg-[#A0BBB1] text-sm">
            Copy
          </button>
        </div>
      </div>
    </div>
  );
};

export default DisplayImageSection;
