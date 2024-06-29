import { useImage } from 'react-image'
import './Book.css'

import { ErrorBoundary } from 'react-error-boundary'

interface BookProps {
    idx: number
    title: string
    img_path: string
    isChosen: boolean
    findSimilarBooks: (id:number, title:string, img_path: string) => void
}

interface BookCoverProps {
    img_path: string
}



const COVER_IMG_PREFIX = 'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/'

const BookCover: React.FC<BookCoverProps> = ({img_path}) =>{

    const {src} = useImage({
        'srcList': COVER_IMG_PREFIX+img_path
    })

    return <img className='book-cover' src={src} alt=''/> 
}


export const Book: React.FC<BookProps> = ({idx, title, img_path, isChosen,findSimilarBooks}) => {

    return (
            <div className={"movie-component " + (isChosen ? "" : "pointer")}
            onClick={() => isChosen ? null : findSimilarBooks(idx, title, img_path)}>
                <ErrorBoundary 
                fallback={<div className='no-image'>
                    {title}
                    </div>}>
                    {img_path != "NaN" ?  <BookCover img_path={img_path}/> : <div>Couldn't load an image</div>}
                </ErrorBoundary>
            </div>
    )
}

