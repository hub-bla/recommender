import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [searchData, setSearchData] = useState("")
  const [requestInputData, setRequestInputData] = useState("")
  useEffect(() => {
    const timeOutId = setTimeout(() => setRequestInputData(searchData), 500)
    return () => clearTimeout(timeOutId)
  }, [searchData])

  useEffect(() => {
    console.log("Tutaj leci request")
    console.log(requestInputData)
  }, [requestInputData])


  
  return (
    <>
      <h1>Movie Recommender</h1>
      <div className="card">
        <input type='text' placeholder='Start searching :)'
        onChange={e => setSearchData(e.target.value)}/>
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>

    </>
  )
}

export default App
