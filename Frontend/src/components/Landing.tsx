import { useEffect, useState } from "react";
import useEmblaCarousel from "embla-carousel-react";
import axios from "axios";
import { Button } from "./ui/button";
export function Landing() {
  const [emblaRef] = useEmblaCarousel();
  const [watches, setWatches] = useState([]);
  useEffect(() => {
    axios
      .get("http://localhost:3000/watches")
      .then((response) => setWatches(response.data))
      .catch((error) => console.error("Error Fetching Data:", error));
  }, []);
  console.log(watches);
  return (
    <div className="embla" ref={emblaRef}>
      <div className="embla__container">
        <div className="embla__slide">
          {watches.map(({ name, price, image }, index) => (
            <div key={index}>
              <Button>Click Me</Button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
