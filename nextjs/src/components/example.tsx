import { StaticImageData } from "next/image";
import React from "react";

interface ExampleProps {
  caption: string;
  image: StaticImageData;
  setDisplay: (image: StaticImageData, caption: string) => void;
}

const Example: React.FC<ExampleProps> = ({ caption, image, setDisplay }) => {
  const handleClick = () => {
    setDisplay(image, caption);
  };
  return (
    <div
      className="h-full w-full cursor-pointer pt-5 shadow-lg"
      onClick={handleClick}
    >
      <img src={image.src} alt={caption} />
    </div>
  );
};

export default Example;
