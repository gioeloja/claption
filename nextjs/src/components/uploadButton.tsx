import React from "react";

interface ButtonProps {
  text: string;
  onClick: () => void;
  disabled?: boolean;
}

const UploadButton: React.FC<ButtonProps> = ({
  text,
  onClick,
  disabled = false,
}) => {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className="h-12 w-36 rounded-lg bg-[#A0BBB1] text-sm hover:bg-[#5D8375] transition"
    >
      {text}
    </button>
  );
};

export default UploadButton;
