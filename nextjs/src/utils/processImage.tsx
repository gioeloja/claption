// Function to convert image file to base64
export function getBase64FromImageFile(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      
      reader.onloadend = () => {
        if (reader.result) {
        // Extract base64 string from data URL
        const base64String = reader.result as string;
        resolve(base64String);
        } else {
          reject(new Error('Failed to convert image to base64'));
        }
      };
      
      reader.onerror = () => {
        reject(new Error('Failed to read the file'));
      };
      
      reader.readAsDataURL(file);
    });
  }
