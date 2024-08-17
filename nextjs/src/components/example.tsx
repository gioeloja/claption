import { StaticImageData } from "next/image";
import React from "react";

interface ExampleProps {
  title: string;
  image: string | StaticImageData;
}

const Example: React.FC<ExampleProps> = ({ title, image }) => {
  return (
    <div className="h-full w-full pt-5 shadow-lg">
      <img src={typeof image === "string" ? image : image.src} alt={title} />
    </div>
  );
};

export default Example;
