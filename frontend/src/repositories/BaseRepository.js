import RepositoryLoggerUtil from '../utils/repositoryLogger'

class BaseRepository {
    constructor(repositoryName) {
        this.repositoryName = repositoryName
        this.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:3160'
    }

    loggedCall(methodName, asyncFunction, params = null) {
        const startTime = RepositoryLoggerUtil.logStart(this.repositoryName, methodName, params)
        
        return asyncFunction()
            .then(result => {
                RepositoryLoggerUtil.logSuccess(this.repositoryName, methodName, startTime, params)
                return result
            })
            .catch(error => {
                RepositoryLoggerUtil.logError(this.repositoryName, methodName, startTime, params)
                throw error
            })
    }

    getHeaders() {
        const token = localStorage.getItem('token')
        return {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    }
}

export default BaseRepository