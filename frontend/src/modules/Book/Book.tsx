import { useImage } from 'react-image'
import './Book.css'

import { ErrorBoundary } from 'react-error-boundary'

interface BookProps {
    idx: number
    title: string
    img_path: string
    findSimilarBooks: (id:number, title:string) => void
}

interface BookCoverProps {
    img_path: string
}



const COVER_IMG_PREFIX = 'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/'

const BookCover: React.FC<BookCoverProps> = ({img_path}) =>{
    // jaka to melodia? thorws error
    if (img_path === "NaN"){
        console.log("hi")
    }
    const {src} = useImage({
        'srcList': COVER_IMG_PREFIX+img_path
    })

    return <img src={src} alt=''/> 
}


export const Book: React.FC<BookProps> = ({idx, title, img_path, findSimilarBooks}) => {

    return (
            <div className='movie-component' onClick={() => findSimilarBooks(idx, title)}>
                <h3>{title}</h3>
                <ErrorBoundary fallback={<div>Couldn't load an image</div>}>
                    {img_path != "NaN" ?  <BookCover img_path={img_path}/> : <div>Couldn't load an image</div>}
                </ErrorBoundary>
              
            </div>
    )
}