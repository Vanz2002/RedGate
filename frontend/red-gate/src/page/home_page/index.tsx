import axios from 'axios';
import React, { useState, useEffect } from 'react';
import { useCookies } from 'react-cookie';

function HomePage() {
  const [plateID, setPlateID] = useState('');
  const [isSubscribe, setIsSubscribe] = useState(false);
  const [plateNumber, setPlateNumber] = useState('');
  const [isRegistering, setIsRegistering] = useState(false);
  const [cookies] = useCookies(["accessToken", "userID"]);

  const instance = axios.create({
    baseURL: "http://127.0.0.1:4444/", // Replace with your API base UR
  });

  useEffect(() => {
    instance
        .get("/plate/getID", {
          headers: {
            Authorization: `Bearer ${cookies.accessToken}`,
          },
        })
        .then((response) => {
          if (response.status == 200) {
            const { v_id, is_subscribe } = response.data as getPlateIDResponseJson;
            setPlateID(v_id); // Set the data in state
            setIsSubscribe(is_subscribe.Bool);
          }
          console.log(response.data);
        })
        .catch((error) => console.error("Error fetching data:", error));
    }, []);

  const handleRegister = async () => {
    const formData = new URLSearchParams();
    formData.append("account_id", cookies.userID);
    formData.append("plate", plateNumber);

    const resp = await instance.post(
        "http://127.0.0.1:4444/plate/create",
        formData,
        {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
        }
      );

    if (resp.data) {
        console.log(resp.data)
    } else {
        console.log("naasjhdhah")
    }
  };

  return (
    <div>
      {isRegistering ? (
        // Step 3: Display a form for registering a plate
        <div>
          <input
            type="text"
            placeholder="Enter Plate Number"
            value={plateNumber}
            onChange={e => setPlateNumber(e.target.value)}
          />
          <button onClick={handleRegister}>Register Plate</button>
        </div>
      ) : (
        // Step 2: Display the plateID if it exists
        <div>
          {plateID ? (
            <p>Plate ID: {plateID}</p>
          ) : (
            // If plateID doesn't exist, display a button to start registration
            <button onClick={() => setIsRegistering(true)}>Register Plate</button>
          )}
        </div>
      )}
    </div>
  );
}

export default HomePage;
