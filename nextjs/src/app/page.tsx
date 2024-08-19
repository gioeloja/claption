"use client";
import Example from "@/components/example";
import UploadFileSection from "@/components/uploadFileSection";
import shirt from "@/assets/shirt.jpeg";
import DisplayImageSection from "@/components/displayImageSection";
import { useReducer } from "react";
import { initialState, PageReducer } from "@/reducer/reducer";
import { StaticImageData } from "next/image";

export default function Home() {
  const [state, dispatch] = useReducer(PageReducer, initialState);

  // image is the base64 encoding data url
  const setDisplay = (image: string, caption: string) => {
    dispatch({
      type: "SET_DISPLAY",
      payload: {
        Image: image,
        Caption: caption,
      },
    });
  };

  const setError = (error: string) => {
    dispatch({
      type: "SET_ERROR",
      payload: {
        Error: error
      },
    });
  }

  const resetDisplay = () => {
    dispatch({ type: "RESET_DISPLAY" });
  };

  const examples = [
    { caption: "example 1", image: "data:image/jpeg;base64,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" },
    { caption: "example 2", image: "data:image/jpeg;base64,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" },
    { caption: "example 3", image: "data:image/jpeg;base64,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" },
  ];

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
          {state.Image && state.Caption ? (
            <DisplayImageSection
              title={state.Caption}
              image={state.Image}
              resetDisplay={resetDisplay}
            />
          ) : (
            <UploadFileSection setDisplay={setDisplay} error={state.Error} setError={setError}/>
          )}
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
          {examples.map((example, index) => (
            <div
              key={index}
              className="h-[250px] w-[250px] md:h-[200px] md:w-[200px]"
            >
              <Example
                caption={example.caption}
                image={example.image}
                setDisplay={setDisplay}
              />
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
