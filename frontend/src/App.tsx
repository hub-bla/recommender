import { Suspense, useEffect, useState } from 'react'
import './App.css'
import { Book } from './modules/Book/Book'
import ReactLoading from 'react-loading';

interface BookInterface {
  id: number
  title: string
  img_path: string
}

type SearchType = "exact" | "semantic"

interface SearchMessage{
  searchbar_input: string
  search_type: SearchType
}

interface RecommendMessage {
  book_id: string
}


function App() {
  const [searchData, setSearchData] = useState("")
  const [requestInputData, setRequestInputData] = useState("")
  const [bookArray, setBookArray] = useState([])
  const [searchType, setSearchType] = useState<SearchType>("exact")
  const [responseFailedError, setResponseFailedError] = useState(false)
  const [notFoundBooksError, setNotFoundBooksError] = useState(false)
  const [similarToTitle, setSimilarToTitle] = useState("")


  const findSimilarBooks = (book_id:number, title: string) => {
    
    const recommendMessage: RecommendMessage = {
      "book_id": book_id.toString()
    }

    const requestOptions = {
      method: 'POST',
      body: JSON.stringify(recommendMessage)
    }

    fetch('http://127.0.0.1:8080/recommend', requestOptions)
          .then(response => response.json())
          .then(data => {
            setBookArray(data['similar_books'] ? data['similar_books'] : [])
            setSimilarToTitle(title)
            setSearchData("")
            setRequestInputData("")
          })
          .catch(() => setResponseFailedError(true))
    

  }


  // block request while user is typing
  useEffect(() => {
    setResponseFailedError(false)


    const timeOutId = setTimeout(() => setRequestInputData(searchData), 500)
    return () => clearTimeout(timeOutId)
  }, [searchData])

  useEffect(() => {
    setResponseFailedError(false)

    
    if (requestInputData != ""){ 

      const searchMessage: SearchMessage = { 
        searchbar_input: requestInputData,
        search_type: searchType
      }

      const requestOptions = {
        method: 'POST',
        body: JSON.stringify(searchMessage)
      }

      fetch('http://127.0.0.1:8080/search', requestOptions)
          .then(response => response.json())
          .then(data => {
            const booksData = data['similar_books'] ? data['similar_books'] : []
            if (booksData.length === 0){
              setNotFoundBooksError(true)
              setBookArray([])
            }else{
              setNotFoundBooksError(false)
              setBookArray(booksData)
            }
          })
          .catch(() => setResponseFailedError(true))
      
      setSimilarToTitle("")
      }
  }, [requestInputData, searchType])


  const books = bookArray.map((book:BookInterface) => <Book 
                                                      key={book.id.toString()} 
                                                      idx={book.id} 
                                                      title={book.title} 
                                                      img_path={book.img_path}
                                                      findSimilarBooks={findSimilarBooks}/>)
  
  console.log(notFoundBooksError)
  return (
    <div className='app-container'>
      <div className='nav-container'>
        <h1>Book Recommender</h1>
          <div className="card">
            {responseFailedError && <p>Error occured</p>}
            <div className='searchbar-container'>
            <input id='searchbar' type='text' placeholder='Start searching :)'
            onChange={e => setSearchData(e.target.value)}  value={searchData}/>
            <div className='search-type-container'>

            <button className={"left-button " + (searchType == "exact" ? "selected-search-type" : "")}
            onClick={() => setSearchType("exact")}>EXACT</button>
            <button className={"right-button " + (searchType == "semantic" ? "selected-search-type" : "")}
            onClick={() => setSearchType("semantic")}>SEMANTIC</button>
            </div>
          </div>
      </div>
      </div>
      <div className='search-results-container'>
          {similarToTitle != "" ? <p>Books similar to: {similarToTitle}</p> : <></>}
          <Suspense fallback={<ReactLoading type={"spinningBubbles"} color={"#213547"} height={'100px'} width={'100px'} />}>
          {notFoundBooksError && <div>Could not find any books</div>}
          {books}
          </Suspense>
        </div>
    </div>
  )
}

export default App
