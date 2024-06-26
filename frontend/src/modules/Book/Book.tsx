import './Book.css'
interface BookProps {
    idx: number
    title: string
    img_path: string
    findSimilarBooks: (id:number, title:string) => void
}

const COVER_IMG_PREFIX = 'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/'

export const Book: React.FC<BookProps> = ({idx, title, img_path, findSimilarBooks}) => {
    return (
        <div className='movie-component' onClick={() => findSimilarBooks(idx, title)}>
            <h3>{title}</h3>
            <img src={COVER_IMG_PREFIX+img_path} alt='' loading='lazy' />
        </div>
    )
}