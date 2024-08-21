import React from "react";
import Image from "next/image";

interface ExampleProps {
  caption: string;
  image: string | string;
  // eslint-disable-next-line no-unused-vars
  setDisplay: (image: string, caption: string) => void;
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
      <Image src={image} alt={caption} />
    </div>
  );
};

export default Example;
