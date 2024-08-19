import React from "react";

interface ExampleProps {
  caption: string;
  image: string | string;
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
      <img src={image} alt={caption} />
    </div>
  );
};

export default Example;
