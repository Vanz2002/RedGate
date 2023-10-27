import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import axios from 'axios'
import LoginForm from './page/login'

function App() {
  const [count, setCount] = useState(0)
  const [payload, setPayload] = useState(null) 

  useEffect(() => {
    // Fetch data from localhost:4444
    axios.get('http://127.0.0.1:4444/')
      .then(response => {
        setPayload(response.data); // Set the data in state
      })
      .catch(error => console.error('Error fetching data:', error));
  }, []); 

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-200">
      <LoginForm />
      <div>
        <h2>Fetched Data:</h2>
        <pre>{JSON.stringify(payload, null, 2)}</pre>
      </div>
    </div>
  );
}

export default App
